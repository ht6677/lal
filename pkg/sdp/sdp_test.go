// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/lal
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package sdp_test

import (
	"testing"

	"github.com/q191201771/lal/pkg/sdp"

	"github.com/q191201771/naza/pkg/nazalog"

	"github.com/q191201771/naza/pkg/assert"
)

var goldenSDP = "v=0" + "\r\n" +
	"o=- 0 0 IN IP6 ::1" + "\r\n" +
	"s=No Name" + "\r\n" +
	"c=IN IP6 ::1" + "\r\n" +
	"t=0 0" + "\r\n" +
	"a=tool:libavformat 57.83.100" + "\r\n" +
	"m=video 0 RTP/AVP 96" + "\r\n" +
	"b=AS:212" + "\r\n" +
	"a=rtpmap:96 H264/90000" + "\r\n" +
	"a=fmtp:96 packetization-mode=1; sprop-parameter-sets=Z2QAIKzZQMApsBEAAAMAAQAAAwAyDxgxlg==,aOvssiw=; profile-level-id=640020" + "\r\n" +
	"a=control:streamid=0" + "\r\n" +
	"m=audio 0 RTP/AVP 97" + "\r\n" +
	"b=AS:30" + "\r\n" +
	"a=rtpmap:97 MPEG4-GENERIC/44100/2" + "\r\n" +
	"a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3; config=1210" + "\r\n" +
	"a=control:streamid=1" + "\r\n"

var goldenSPS = []byte{
	0x67, 0x64, 0x00, 0x20, 0xAC, 0xD9, 0x40, 0xC0, 0x29, 0xB0, 0x11, 0x00, 0x00, 0x03, 0x00, 0x01, 0x00, 0x00, 0x03, 0x00, 0x32, 0x0F, 0x18, 0x31, 0x96,
}

var goldenPPS = []byte{
	0x68, 0xEB, 0xEC, 0xB2, 0x2C,
}

func TestParseSDP(t *testing.T) {
	sdpCtx, err := sdp.ParseSDP([]byte(goldenSDP))
	assert.Equal(t, nil, err)
	nazalog.Debugf("sdp=%+v", sdpCtx)
}

func TestParseARTPMap(t *testing.T) {
	golden := map[string]sdp.ARTPMap{
		"rtpmap:96 H264/90000": {
			PayloadType:        96,
			EncodingName:       "H264",
			ClockRate:          90000,
			EncodingParameters: "",
		},
		"rtpmap:97 MPEG4-GENERIC/44100/2": {
			PayloadType:        97,
			EncodingName:       "MPEG4-GENERIC",
			ClockRate:          44100,
			EncodingParameters: "2",
		},
	}
	for in, out := range golden {
		actual, err := sdp.ParseARTPMap(in)
		assert.Equal(t, nil, err)
		assert.Equal(t, out, actual)
	}
}

func TestParseFmtPBase(t *testing.T) {
	golden := map[string]sdp.AFmtPBase{
		"a=fmtp:96 packetization-mode=1; sprop-parameter-sets=Z2QAIKzZQMApsBEAAAMAAQAAAwAyDxgxlg==,aOvssiw=; profile-level-id=640020": {
			Format: 96,
			Parameters: map[string]string{
				"packetization-mode":   "1",
				"sprop-parameter-sets": "Z2QAIKzZQMApsBEAAAMAAQAAAwAyDxgxlg==,aOvssiw=",
				"profile-level-id":     "640020",
			},
		},
		"a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3; config=1210": {
			Format: 97,
			Parameters: map[string]string{
				"profile-level-id": "1",
				"mode":             "AAC-hbr",
				"sizelength":       "13",
				"indexlength":      "3",
				"indexdeltalength": "3",
				"config":           "1210",
			},
		},
	}
	for in, out := range golden {
		actual, err := sdp.ParseAFmtPBase(in)
		assert.Equal(t, nil, err)
		assert.Equal(t, out, actual)
	}
}

func TestParseSPSPPS(t *testing.T) {
	s := "a=fmtp:96 packetization-mode=1; sprop-parameter-sets=Z2QAIKzZQMApsBEAAAMAAQAAAwAyDxgxlg==,aOvssiw=; profile-level-id=640020"
	f, err := sdp.ParseAFmtPBase(s)
	assert.Equal(t, nil, err)
	sps, pps, err := sdp.ParseSPSPPS(f)
	assert.Equal(t, nil, err)
	assert.Equal(t, goldenSPS, sps)
	assert.Equal(t, goldenPPS, pps)
}

func TestParseASC(t *testing.T) {
	s := "a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3; config=1210"
	f, err := sdp.ParseAFmtPBase(s)
	assert.Equal(t, nil, err)
	asc, err := sdp.ParseASC(f)
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte{0x12, 0x10}, asc)
}
