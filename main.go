package main

import (
	"fmt"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// CubicCircuit defines a simple circuit
// x**3 + x + 5 == y
type CubicCircuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
// x**3 + x + 5 == y
func (circuit *CubicCircuit) Define(api frontend.API) error {
	x3 := api.Mul(circuit.X, circuit.X, circuit.X)
	for i := 0; i < 1000000; i++ {
		api.AssertIsEqual(circuit.Y, api.Add(x3, circuit.X, 5))
	}
	return nil
}

func main() {
	// compiles our circuit into a R1CS
	var circuit CubicCircuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(ccs)

	// witness definition
	assignment := CubicCircuit{X: 3, Y: 35}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()
	// groth16: Prove & Verify
	icicleTimeStart := time.Now()
	proof, err := groth16.Prove(ccs, pk, witness, backend.WithIcicleAcceleration())
	fmt.Println("Icicle proof time:", time.Since(icicleTimeStart))

	if err != nil {
		fmt.Println(err)
	}

	if err = groth16.Verify(proof, vk, publicWitness); err != nil {
		panic("Failed verification")
	}
	
	gnarkTimeStart := time.Now()
	proofNoIcicle, err := groth16.Prove(ccs, pk, witness)
	fmt.Println("Gnark CPU proof time:", time.Since(gnarkTimeStart))
	
	if err != nil {
		fmt.Println(err)
	}

	if err = groth16.Verify(proofNoIcicle, vk, publicWitness); err != nil {
		panic("Failed verification")
	}
}
