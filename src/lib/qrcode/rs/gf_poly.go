package rs

import (
	"log"

	"yotudo/src/lib/qrcode/bitbuffer"
)

type gfPoly []gfElement

func newGFPolyFromData(data *bitbuffer.BitBuffer) gfPoly {
	numTotalBytes := data.Len() / 8
	if data.Len()%8 != 0 {
		numTotalBytes++
	}

	result := make(gfPoly, numTotalBytes)

	for i, j := numTotalBytes-1, uint(0); j < data.Len(); j += 8 {
		result[i] = gfElement(data.ByteAt(j))
		i--
	}

	return result
}

func newGFPolyMonomial(term gfElement, degree int) gfPoly {
	if term == gfZero {
		return gfPoly{}
	}

	result := make(gfPoly, degree+1)
	result[degree] = term

	return result
}

func (e gfPoly) data(numTerms int) []byte {
	result := make([]byte, numTerms)

	for i, j := numTerms-len(e), len(e)-1; j >= 0; j-- {
		result[i] = byte(e[j])
		i++
	}

	return result
}

func gfPolyMultiply(a, b gfPoly) gfPoly {
	numATerms := len(a)
	numBTerms := len(b)

	result := make(gfPoly, numATerms+numBTerms)

	for i := 0; i < numATerms; i++ {
		for j := 0; j < numBTerms; j++ {
			if a[i] != 0 && b[j] != 0 {
				monomial := make(gfPoly, i+j+1)
				monomial[i+j] = gfMultiply(a[i], b[j])

				result = gfPolyAdd(result, monomial)
			}
		}
	}

	return result.normalised()
}

func gfPolyRemainder(numerator, denominator gfPoly) gfPoly {
	if denominator.equals(gfPoly{}) {
		log.Panicln("Remainder by zero")
	}

	remainder := numerator

	for len(remainder) >= len(denominator) {
		degree := len(remainder) - len(denominator)
		coefficient := gfDivide(remainder[len(remainder)-1], denominator[len(denominator)-1])
		divisor := gfPolyMultiply(denominator, newGFPolyMonomial(coefficient, degree))
		remainder = gfPolyAdd(remainder, divisor)
	}

	return remainder.normalised()
}

func gfPolyAdd(a, b gfPoly) gfPoly {
	numATerms := len(a)
	numBTerms := len(b)

	numTerms := numATerms
	if numBTerms > numTerms {
		numTerms = numBTerms
	}

	result := make(gfPoly, numTerms)

	for i := 0; i < numTerms; i++ {
		switch {
		case numATerms > i && numBTerms > i:
			result[i] = gfAdd(a[i], b[i])
		case numATerms > i:
			result[i] = a[i]
		default:
			result[i] = b[i]
		}
	}

	return result.normalised()
}

func (e *gfPoly) normalised() gfPoly {
	numTerms := len(*e)
	maxNonzeroTerm := numTerms - 1

	for i := numTerms - 1; i >= 0; i-- {
		if (*e)[i] != 0 {
			break
		}

		maxNonzeroTerm = i - 1
	}

	if maxNonzeroTerm < 0 {
		return gfPoly{}
	} else if maxNonzeroTerm < numTerms-1 {
		*e = (*e)[0 : maxNonzeroTerm+1]
	}

	return *e
}

func (e gfPoly) equals(other gfPoly) bool {
	var minecPoly *gfPoly
	var maxecPoly *gfPoly

	if len(e) > len(other) {
		minecPoly = &other
		maxecPoly = &e
	} else {
		minecPoly = &e
		maxecPoly = &other
	}

	numMinTerms := len(*minecPoly)
	numMaxTerms := len(*maxecPoly)

	for i := 0; i < numMinTerms; i++ {
		if e[i] != other[i] {
			return false
		}
	}

	for i := numMinTerms; i < numMaxTerms; i++ {
		if (*maxecPoly)[i] != 0 {
			return false
		}
	}

	return true
}
