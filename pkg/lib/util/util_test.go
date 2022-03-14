package util

import (
	"testing"
)

func TestDecimalToRaw(t *testing.T) {
	type args struct {
		value   string
		decimal int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{value: "0.00002000", decimal: 8},
			want: "2000",
		},
		{
			name: "ok",
			args: args{value: "96092252354.64214000", decimal: 8},
			want: "9609225235464214000",
		},
		{
			name: "ok",
			args: args{value: "100092252354.64214000", decimal: 8},
			want: "10009225235464214000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecimalToRaw(tt.args.value, tt.args.decimal)
			// log.Println("masok3", err)
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("RawToCoin() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			if got != tt.want {
				t.Errorf("RawToCoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawToDecimal(t *testing.T) {
	type args struct {
		value   string
		decimal int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{value: "2000", decimal: 8},
			want: "0.00002",
		},
		{
			name: "ok",
			args: args{value: "9609225235464214001", decimal: 8},
			want: "96092252354.64214001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RawToDecimal(tt.args.value, tt.args.decimal)
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("RawToCoin() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			if got != tt.want {
				t.Errorf("RawToCoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestFormatCurrency(t *testing.T) {
// 	type args struct {
// 		value string
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		wantResult string
// 	}{
// 		{
// 			args:       args{value: "12345678"},
// 			wantResult: "12,345,678",
// 		},
// 		{
// 			args:       args{value: "1234567"},
// 			wantResult: "1,234,567",
// 		},
// 		{
// 			args:       args{value: "123456"},
// 			wantResult: "123,456",
// 		},
// 		{
// 			args:       args{value: "1234"},
// 			wantResult: "1,234",
// 		},
// 		{
// 			args:       args{value: "1234.222"},
// 			wantResult: "1,234.222",
// 		},
// 		{
// 			args:       args{value: "96092252354.64214000"},
// 			wantResult: "96,092,252,354.64214000",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if gotResult := FormatCurrency(tt.args.value); gotResult != tt.wantResult {
// 				t.Errorf("FormatCurrency() = %v, want %v", gotResult, tt.wantResult)
// 			}
// 		})
// 	}
// }
