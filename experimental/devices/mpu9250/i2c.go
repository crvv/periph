package mpu9250

import (
	"fmt"

	"periph.io/x/periph/conn/i2c"
)

type I2CTransport struct {
	dev i2c.Dev
}

func NewI2CTransport(b i2c.Bus, addr uint16) (*I2CTransport, error) {
	t := I2CTransport{dev: i2c.Dev{Bus: b, Addr: addr}}
	return &t, nil
}

func (s *I2CTransport) writeByte(address byte, value byte) error {
	_, err := s.dev.Write([]byte{address, value})
	return err
}

func (s *I2CTransport) writeMagReg(address byte, value byte) error {
	return s.writeByte(address, value)
}

func (s *I2CTransport) writeMaskedReg(address byte, mask byte, value byte) error {
	maskedValue := mask & value
	regVal, err := s.readByte(address)
	if err != nil {
		return err
	}
	regVal = (regVal &^ maskedValue) | maskedValue
	return s.writeByte(address, regVal)
}

func (s *I2CTransport) readMaskedReg(address byte, mask byte) (byte, error) {
	reg, err := s.readByte(address)
	if err != nil {
		return 0, err
	}
	return reg & mask, nil
}

func (s *I2CTransport) readByte(address byte) (byte, error) {
	var res [1]byte
	err := s.dev.Tx([]byte{address}, res[:])
	return res[0], err
}

func (s *I2CTransport) readUint16(address ...byte) (uint16, error) {
	if len(address) != 2 {
		return 0, fmt.Errorf("only 2 bytes per read")
	}
	h, err := s.readByte(address[0])
	if err != nil {
		return 0, err
	}
	l, err := s.readByte(address[1])
	if err != nil {
		return 0, err
	}
	return uint16(h)<<8 | uint16(l), nil
}

var _ Proto = &I2CTransport{}
