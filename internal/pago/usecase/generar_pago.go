package usecase

import (
	"opennode/internal/pago/model"
	"opennode/internal/pago/port"
)

type PagoUsecase struct {
	pagoPort port.PagoPort
}

func NewPagoUsecase(p port.PagoPort) PagoUsecase {
	return PagoUsecase{
		pagoPort: p,
	}
}

func (p PagoUsecase) RealizarPago(pago model.Pago) (model.Pago, error) {
	return p.pagoPort.RealizarPago(pago)
}
