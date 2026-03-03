package lib

func Dividir(numerador, denominador float64) float64 {
	if denominador == 0 {
		return 0
	}
	return numerador / denominador
}
