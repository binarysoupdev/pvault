package meta

import "io"

type EncoderMock struct {
	Metadata Metadata

	EncodeMetadataError error
	DecodeMetadataError error
}

func (mock *EncoderMock) EncodeMetadata(_ io.Writer, m Metadata) error {
	mock.Metadata = m
	return mock.EncodeMetadataError
}

func (m EncoderMock) DecodeMetadata(_ io.Reader) (Metadata, error) {
	return m.Metadata, m.DecodeMetadataError
}
