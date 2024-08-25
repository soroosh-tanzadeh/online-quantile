# Online Quantile Estimation
An efficient implementation of P-Squred algorithm for online quantile estimation, allowing you to estimate quantiles (such as medians or percentiles) in a data stream without the need to store all data points.

**Source Research Paper**: "The P2 Algorithm for Dynamic Calculation of Quantiles and Histograms Without Storing Observations"; Raj Jain, Imrich Chlamtac

[Link to the paper](https://www.cse.wustl.edu/~jain/papers/ftp/psqr.pdf)

## Example
```Go
package main

import (
	"fmt"
	onlinequantile "github.com/soroosh-tanzadeh/online-quantile"
)

func main() {
	series := []float64{0.02, 0.5, 0.74, 3.39, 0.83, 22.37, 10.15, 15.43, 38.62, 15.92, 34.60, 10.28, 1.47, 0.40, 0.05, 1.39, 0.27, 0.42, 0.09, 11.37}

	q50 := onlinequantile.NewQuantile(0.5)

	// Initialize with first 5 points
	q50.Update(series[:5])

	for _, x := range series[5:] {
		q50.Consume(x)
		fmt.Printf("P50= %f \n", q50.GetQuantileValue())
	}
}
```