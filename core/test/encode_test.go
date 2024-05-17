package test

import (
	"CollabEdit/core"
	"bytes"
	"testing"
)

// TestEncoderDecoder 正常编码和解码测试
func TestEncoderDecoder(t *testing.T) {
	encoder := core.NewEncoder()
	encoder.WriteVarUint(89460623423423)
	encoder.WriteString("Hello World 12345678!")
	buf := encoder.Bytes()

	// 解码
	decoder := core.NewDecoder(buf)

	// 读取VarUint
	varUint, err := decoder.ReadVarUint()
	if err != nil {
		t.Fatalf("读取VarUint时出错: %v", err)
	}
	expectedVarUint := uint64(89460623423423)
	if varUint != expectedVarUint {
		t.Fatalf("期望 %d，得到 %d", expectedVarUint, varUint)
	}

	// 检查是否还有内容未读取
	if !decoder.HasContent() {
		t.Fatal("期望还有内容未读取")
	}

	// 读取字符串
	readString, err := decoder.ReadString()
	if err != nil {
		t.Fatalf("读取字符串时出错: %v", err)
	}
	expectedString := "Hello World 12345678!"
	if readString != expectedString {
		t.Fatalf("期望 %s，得到 %s", expectedString, readString)
	}

	// 检查是否还有内容未读取
	if decoder.HasContent() {
		t.Fatal("期望没有剩余内容")
	}
}

// TestDecoderWithTruncatedData 测试解码器处理截断数据
func TestDecoderWithTruncatedData(t *testing.T) {
	data := []byte{0xA3, 0x02} // 不完整的VarUint
	decoder := core.NewDecoder(data)
	_, err := decoder.ReadVarUint()
	if err != nil {
		t.Fatal("期望因截断数据而出错:", err)
	}
}

// TestDecoderWithInvalidStringLength 测试解码器处理无效的字符串长度
func TestDecoderWithInvalidStringLength(t *testing.T) {
	encoder := core.NewEncoder()
	encoder.WriteVarUint(100) // 假设字符串长度为100，但实际数据不足
	buf := encoder.Bytes()

	decoder := core.NewDecoder(buf)
	_, err := decoder.ReadString()
	if err == nil {
		t.Fatal("期望因无效的字符串长度而出错，得到nil")
	}
}

// TestDecoderWithExcessiveVarUintValue 测试解码器处理过大的VarUint值
func TestDecoderWithExcessiveVarUintValue(t *testing.T) {
	// 创建一个超过64位限制的VarUint
	data := []byte{0x81, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x02}
	decoder := core.NewDecoder(data)
	_, err := decoder.ReadVarUint()
	if err != nil {
		t.Fatal("期望因过大的VarUint值而出错，得到nil")
	}
}

// TestDecoderWithExtraContent 测试解码器处理额外内容
func TestDecoderWithExtraContent(t *testing.T) {
	encoder := core.NewEncoder()
	encoder.WriteVarUint(1)
	encoder.WriteString("Test")
	encoder.Write([]byte{0x01, 0x02, 0x03}) // 添加额外内容
	buf := encoder.Bytes()

	decoder := core.NewDecoder(buf)
	_, err := decoder.ReadVarUint()
	if err != nil {
		t.Fatalf("读取VarUint时出错: %v", err)
	}
	_, err = decoder.ReadString()
	if err != nil {
		t.Fatalf("读取字符串时出错: %v", err)
	}

	if !decoder.HasContent() {
		t.Fatal("期望还有内容未读取")
	}

	// 读取额外内容
	extraContent := make([]byte, 3)
	_, err = decoder.Read(extraContent)
	if err != nil {
		t.Fatalf("读取额外内容时出错: %v", err)
	}
	expectedExtraContent := []byte{0x01, 0x02, 0x03}
	if !bytes.Equal(extraContent, expectedExtraContent) {
		t.Fatalf("期望额外内容 %v，得到 %v", expectedExtraContent, extraContent)
	}
}
