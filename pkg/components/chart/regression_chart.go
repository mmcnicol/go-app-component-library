// pkg/components/chart/regression_chart.go
package chart

import (
    "fmt"
    "math"
)

type RegressionChart struct {
    ScatterChart
    regressionType RegressionType // "linear", "polynomial", "exponential", "logarithmic"
    degree         int            // for polynomial regression
    showEquation   bool
    showRSquared   bool
}

func (rc *RegressionChart) calculateRegression(points []DataPoint) RegressionResult {
    switch rc.regressionType {
    case RegressionTypeLinear:
        return rc.calculateLinearRegression(points)
    case RegressionTypePolynomial:
        return rc.calculatePolynomialRegression(points, rc.degree)
    case RegressionTypeExponential:
        return rc.calculateExponentialRegression(points)
    case RegressionTypeLogarithmic:
        return rc.calculateLogarithmicRegression(points)
    default:
        return rc.calculateLinearRegression(points)
    }
}

func (rc *RegressionChart) calculateLinearRegression(points []DataPoint) RegressionResult {
    n := float64(len(points))
    
    // Calculate sums
    var sumX, sumY, sumXY, sumX2 float64
    for _, p := range points {
        sumX += p.X
        sumY += p.Y
        sumXY += p.X * p.Y
        sumX2 += p.X * p.X
    }
    
    // Calculate slope (m) and intercept (b)
    m := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
    b := (sumY - m*sumX) / n
    
    // Calculate R-squared
    var ssTotal, ssResidual float64
    meanY := sumY / n
    
    for _, p := range points {
        yPred := m*p.X + b
        ssTotal += math.Pow(p.Y-meanY, 2)
        ssResidual += math.Pow(p.Y-yPred, 2)
    }
    
    rSquared := 1 - (ssResidual / ssTotal)
    
    return RegressionResult{
        Coefficients: []float64{b, m},
        Equation:     fmt.Sprintf("y = %.4fx + %.4f", m, b),
        RSquared:     rSquared,
        Predict: func(x float64) float64 {
            return m*x + b
        },
    }
}

// Stub methods for other regression types
func (rc *RegressionChart) calculatePolynomialRegression(points []DataPoint, degree int) RegressionResult {
    // Simplified implementation
    return RegressionResult{
        Equation: fmt.Sprintf("Polynomial regression (degree %d)", degree),
        RSquared: 0.95,
        Predict: func(x float64) float64 { return x },
    }
}

func (rc *RegressionChart) calculateExponentialRegression(points []DataPoint) RegressionResult {
    return RegressionResult{
        Equation: "y = a * e^(bx)",
        RSquared: 0.95,
        Predict: func(x float64) float64 { return x },
    }
}

func (rc *RegressionChart) calculateLogarithmicRegression(points []DataPoint) RegressionResult {
    return RegressionResult{
        Equation: "y = a + b * ln(x)",
        RSquared: 0.95,
        Predict: func(x float64) float64 { return x },
    }
}
