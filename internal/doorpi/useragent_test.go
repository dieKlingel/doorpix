package doorpi_test

import (
	"github.com/dieklingel/doorpix/internal/transport/sip"
	"github.com/stretchr/testify/mock"
)

type MockUserAgent struct {
	mock.Mock
}

func (m *MockUserAgent) Invite(uri string) (*sip.CallInfo, error) {
	args := m.Called(uri)
	return args.Get(0).(*sip.CallInfo), args.Error(1)
}

func (m *MockUserAgent) SendMessage(uri string, body string) error {
	args := m.Called(uri, body)
	return args.Error(0)
}
