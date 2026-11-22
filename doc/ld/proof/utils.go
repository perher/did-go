/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package proof

import (
	"errors"
)

const (
	jsonldProof = "proof"
)

// GetProofs gets proof(s) from LD Object.
func GetProofs(jsonLdObject map[string]any) ([]*Proof, error) {
	entry, ok := jsonLdObject[jsonldProof]
	if !ok {
		return nil, ErrProofNotFound
	}

	var typedEntry []any
	switch te := entry.(type) {
	case []any:
		typedEntry = te
	case map[string]any:
		typedEntry = []any{te}
	default:
		return nil, errors.New("expecting []interface{} or map[string]interface{}, got something else")
	}

	var result []*Proof

	for _, e := range typedEntry {
		emap, ok := e.(map[string]any)
		if !ok {
			return nil, errors.New("wrong interface, expecting []interface{}")
		}

		proof, err := NewProof(emap)
		if err != nil {
			return nil, err
		}

		result = append(result, proof)
	}

	return result, nil
}

// AddProof adds a proof to LD Object.
func AddProof(jsonLdObject map[string]any, proof *Proof) error {
	var proofs []any

	entry, exists := jsonLdObject[jsonldProof]

	if exists {
		switch p := entry.(type) {
		case []any:
			proofs = p
		default:
			proofs = []any{p}
		}
	}

	proofs = append(proofs, proof.JSONLdObject())
	jsonLdObject[jsonldProof] = proofs

	return nil
}

// GetCopyWithoutProof gets copy of JSON LD Object without proofs (signatures).
func GetCopyWithoutProof(jsonLdObject map[string]any) map[string]any {
	if jsonLdObject == nil {
		return nil
	}

	dest := make(map[string]any)

	for k, v := range jsonLdObject {
		if k != jsonldProof {
			dest[k] = v
		}
	}

	return dest
}

// ErrProofNotFound is returned when proof is not found.
var ErrProofNotFound = errors.New("proof not found")
