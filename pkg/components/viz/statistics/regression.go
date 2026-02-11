// pkg/components/viz/statistics/regression.go
package statistics

import "math"

// Regression types
type RegressionType string
const (
    Linear      RegressionType = "linear"
    Polynomial  RegressionType = "polynomial"
    Exponential RegressionType = "exponential"
    Logarithmic RegressionType = "logarithmic"
    Power       RegressionType = "power"
)

// RegressionResult contains regression analysis results
type RegressionResult struct {
    Type        RegressionType
    Coefficients []float64
    Equation    string
    RSquared    float64
    RMSE        float64
    Predict     func(x float64) float64
}

// CalculateRegression performs regression analysis on data points
func CalculateRegression(points []viz.Point, rtype RegressionType, degree int) RegressionResult {
    switch rtype {
    case Linear:
        return calculateLinearRegression(points)
    case Polynomial:
        return calculatePolynomialRegression(points, degree)
    case Exponential:
        return calculateExponentialRegression(points)
    case Logarithmic:
        return calculateLogarithmicRegression(points)
    case Power:
        return calculatePowerRegression(points)
    default:
        return calculateLinearRegression(points)
    }
}
