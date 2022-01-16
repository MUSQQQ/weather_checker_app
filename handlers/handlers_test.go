package handlers

import (
	"testing"
	"weather_checker/models"

	"github.com/stretchr/testify/require"
)

func TestByteArrayToFloat(t *testing.T) {
	tests := []struct {
		input      []byte
		wantOutput float32
		wantErr    bool
	}{
		{
			[]byte("123.321"),
			123.321,
			false,
		},
		{
			[]byte("64.64"),
			64.64,
			false,
		},
		{
			[]byte("100"),
			100.0,
			false,
		},
		{
			[]byte("A#%)(+"),
			0,
			true,
		},
	}

	for _, tt := range tests {
		got, err := byteArrayToFloat(tt.input)
		if tt.wantErr {
			require.Error(t, err)
		} else {
			require.Equal(t, tt.wantOutput, got)
		}
	}
}

func TestDetermineOverall(t *testing.T) {
	tests := []struct {
		input      models.Weather
		wantOutput string
	}{
		{
			models.Weather{
				Temperature: 60,
				Pressure:    1013,
				Humidity:    80,
				Clouds:      20,
			},
			"Hot and sunny",
		},
		{
			models.Weather{
				Temperature: -10,
				Pressure:    1013,
				Humidity:    0,
				Clouds:      50,
			},
			"Cold and cloudy",
		},
		{
			models.Weather{
				Temperature: 15,
				Pressure:    1013,
				Humidity:    80,
				Clouds:      80,
			},
			"Warm and cloudy",
		},
	}
	for _, tt := range tests {
		got := determineOverall(tt.input)
		require.Equal(t, tt.wantOutput, got)
	}
}
