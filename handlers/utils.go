package handlers

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
	"weather_checker/models"
)

func determineOverall(data *models.Weather) (overall string) {
	switch {
	case data.Temperature < 0.0:
		{
			switch {
			case data.Clouds < 50.0:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Cold, possible rain"
					default:
						overall = "Cold and the sky is clear"
					}

				}
			default:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Cold, cloudy and possible rain"
					default:
						overall = "Cold and cloudy"
					}

				}

			}

		}
	case data.Temperature < 10.0:
		{
			switch {
			case data.Clouds < 50.0:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Cool, possible rain"
					default:
						overall = "Cool and the sky is clear"
					}

				}
			default:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Cool, cloudy and possible rain"
					default:
						overall = "Cool and cloudy"
					}

				}

			}
		}
	case data.Temperature < 20.0:
		{
			switch {
			case data.Clouds < 50.0:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Warm, possible rain"
					default:
						overall = "Warm and the sky is clear"
					}

				}
			default:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Warm, cloudy and possible rain"
					default:
						overall = "Warm and cloudy"
					}

				}

			}

		}
	default:
		{
			switch {
			case data.Clouds < 50.0:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Hot, possible rain"
					default:
						overall = "Hot and the sky is clear"
					}

				}
			default:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Hot, cloudy and possible rain"
					default:
						overall = "Hot and cloudy"
					}

				}

			}
		}
	}
	return overall
}

func timestampToString(timestamp, offset int64) (result string) {
	year, month, day := time.Unix(timestamp, 0).UTC().Add(time.Duration(offset) * time.Second).Date()
	hour, minute, _ := time.Unix(timestamp, 0).UTC().Add(time.Duration(offset) * time.Second).Clock()

	possibleOneDigit := []int{int(month), day, hour, minute}
	finalDoubleDigitResult := []string{"", "", "", ""}
	for i := 0; i < 4; i++ {
		if possibleOneDigit[i]/10 == 0 {
			finalDoubleDigitResult[i] = fmt.Sprintf("0%d", possibleOneDigit[i])
		} else {
			finalDoubleDigitResult[i] = fmt.Sprintf("%d", possibleOneDigit[i])
		}
	}
	result = fmt.Sprintf("%d:%s:%s\n%s:%s", year, finalDoubleDigitResult[0], finalDoubleDigitResult[1], finalDoubleDigitResult[2], finalDoubleDigitResult[3])
	return result
}

func byteArrayToFloat(bytes []byte) (result float32, err error) {
	strByte := string(bytes)
	fmt.Println(strByte)
	var i int
	for i = 0; i < len(strByte); i++ {
		if strByte[i] == '.' {
			break
		}
	}

	var intResultPart int

	intResultPart, err = strconv.Atoi(strByte[0:i])

	if err != nil {
		log.Printf("error while converting: %s", err)
		return -1, err
	}

	var mantissaPart int
	if i >= len(strByte)-1 {
		mantissaPart = 0

	} else {
		mantissaPart, err = strconv.Atoi(strByte[i+1:])
	}

	if err != nil {
		log.Printf("error while converting: %s", err)
		return -1, err
	}
	result = float32(intResultPart) + float32(mantissaPart)/(float32(math.Pow10(len(strByte)-i-1)))
	return result, nil
}
