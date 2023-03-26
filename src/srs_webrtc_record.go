package main

import (
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"os"
	"strings"
	"webrtc-exmple/config"
	"webrtc-exmple/entity/srs"
	"webrtc-exmple/internal/signal"
	"webrtc-exmple/utils"

	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v3"
)

func init() {
	config.InitLocalLog()
}

var play_stream = "webrtc://139.159.213.37:10985/live/livestream"
var srs_api = "http://139.159.213.37:10985/rtc/v1/play/"

func main() {

	// Create a MediaEngine object to configure the supported codec
	m := &webrtc.MediaEngine{}

	// Setup the codecs you want to use.
	// We'll use a VP8 and Opus but you can also define your own
	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264, ClockRate: 90000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        127,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		seelog.Errorf("register video codec error:%v", err)
		panic(err)
	}
	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		seelog.Errorf("register audio codec error:%v", err)
		panic(err)
	}

	i := &interceptor.Registry{}

	// Use the default set of Interceptors
	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		seelog.Errorf("register default interceptor error:%v", err)
		panic(err)
	}

	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i))

	// Prepare the configuration
	config := webrtc.Configuration{}

	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		seelog.Errorf("api new peer connection error:%v", err)
		panic(err)
	}

	// Allow us to receive 1 audio track, and 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		seelog.Errorf("peer connection add audio transceiver error:%v", err)
		panic(err)
	}
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		seelog.Errorf("peer connection add video transceiver error:%v", err)
		panic(err)
	}

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		seelog.Errorf("peer connection create offer err:%v", err)
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		seelog.Errorf("peer connection set local sdp error:%v", err)
		panic(err)
	}

	// Set a handler for when a new remote track starts, this handler saves buffers to disk as
	// an ivf file, since we could have multiple video tracks we provide a counter.
	// In your application this is where you would handle/process video
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		codec := track.Codec()
		if strings.EqualFold(codec.MimeType, webrtc.MimeTypeOpus) {
			seelog.Infof("Got Opus track, saving to disk as output.opus (48 kHz, 2 channels)")
		} else if strings.EqualFold(codec.MimeType, webrtc.MimeTypeH264) {
			seelog.Infof("Got H264 track, saving to disk as output.ivf")
			recordFrame(track)
		}
	})

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		seelog.Infof("Connection State has changed %s", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateConnected {
			seelog.Info("Ctrl+C the remote client to stop the demo")
		} else if connectionState == webrtc.ICEConnectionStateFailed {
			seelog.Errorf("peer connection stata failure")

			// Gracefully shutdown the peer connection
			if closeErr := peerConnection.Close(); closeErr != nil {
				seelog.Errorf("peer connection close error:%v", closeErr)
				panic(closeErr)
			}

			os.Exit(0)
		}
	})
	answerSdp := srsCall(offer.SDP)

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answerSdp,
	})
	if err != nil {
		seelog.Errorf("peer connection set remote sdp %s error:%v", answerSdp, err)
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(signal.Encode(*peerConnection.LocalDescription()))
	seelog.Infof("peer connection local sdp:%s", signal.Encode(*peerConnection.LocalDescription()))

	// Block forever
	select {}
}

func srsCall(sdp string) string {
	req := &srs.SdpRequest{
		Api:       srs_api,
		ClientIp:  "",
		StreamURL: play_stream,
		Sdp:       sdp,
	}
	body, _ := json.Marshal(req)
	//seelog.Infof("srs server call, body:%s", string(body))

	if data, err := utils.HttpPost(srs_api, body); err == nil {
		var resp srs.SdpResponse
		if json.Unmarshal(data, &resp) == nil {
			if resp.Code == 0 {
				return resp.Sdp
			}
			seelog.Errorf("srs call failure. resp:%s", string(data))
		}
	}
	return ""
}

func recordFrame(track *webrtc.TrackRemote) {
	index := 0
	for {
		rtpPacket, _, err := track.ReadRTP()
		if err != nil {
			panic(err)
		}

		seelog.Infof("receive rtp packet index=%d, Payload length:%v", index, len(rtpPacket.Payload))
		index++

	}
}
