package rs

import (
	"log"

	"yotudo/src/lib/qrcode/bitbuffer"
)

func Encode(data *bitbuffer.BitBuffer, numECBytes int) *bitbuffer.BitBuffer {
	ecpoly := newGFPolyFromData(data)
	ecpoly = gfPolyMultiply(ecpoly, newGFPolyMonomial(gfOne, numECBytes))

	generator := rsGeneratorPoly(numECBytes)

	remainder := gfPolyRemainder(ecpoly, generator)

	result := bitbuffer.Clone(data)
	result.AppendBytes(remainder.data(numECBytes))

	return result
}

func rsGeneratorPoly(degree int) gfPoly {
	if degree < 2 {
		log.Panic("degree less than 2")
	}

	var generator gfPoly = []gfElement{1}

	for i := 0; i < degree; i++ {
		var nextPoly gfPoly = []gfElement{gfExpTable[i], 1}
		generator = gfPolyMultiply(generator, nextPoly)
	}

	return generator
}
