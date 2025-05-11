package blockchain

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Verify wallet_signature
func VerifyWalletSignature(walletAddress, message, signature string) (bool, error) {
	addr := common.HexToAddress(walletAddress)

	// 2. Create hashPrefecture
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))
	data := []byte(prefix + message)
	hash := crypto.Keccak256Hash(data)

	// 3.  (bỏ "0x" nếu có)
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature format: %v", err)
	}

	// 4. Kiểm tra xem chữ ký có đúng định dạng không
	if len(sig) != 65 {
		return false, fmt.Errorf("invalid signature length")
	}

	// 5. Điều chỉnh V trong chữ ký (EIP-155)
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	// 6. Khôi phục khóa công khai từ chữ ký
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return false, fmt.Errorf("error recovering public key: %v", err)
	}

	// 7. Lấy địa chỉ từ khóa công khai
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// 8. So sánh với địa chỉ ví được cung cấp
	return addr == recoveredAddr, nil
}

// Hàm lưu dữ liệu KYC
func SaveKYCData(citizen_id, wallet_address, full_name, phone_number, date_of_birth, nationality, wallet_signature string) error {
	message := fmt.Sprintf("I confirm that I am the owner of wallet address %s and the information provided for KYC is accurate.", wallet_address)

	isValid, err := VerifyWalletSignature(wallet_address, message, wallet_signature)
	if err != nil {
		return fmt.Errorf("signature verification error: %v", err)
	}

	if !isValid {
		return fmt.Errorf("signature verification failed: wallet signature does not match wallet address")
	}

	// 3. Nếu chữ ký hợp lệ, lưu dữ liệu vào database
	// db.Exec("INSERT INTO KYC (...) VALUES (...)")

	return nil
}