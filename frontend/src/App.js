import React, { useState, useCallback } from 'react';
import './App.css';

const MatrixDisplay = React.memo(({ matrixInputs, onInputChange, title, editable = false }) => (
  <div className="matrix-container">
    <h3>{title}</h3>
    <div className="matrix">
      {matrixInputs.map((row, rowIdx) => (
        <div key={rowIdx} className="matrix-row">
          <span className="matrix-bracket">[</span>
          {row.map((cell, colIdx) => (
            <input
              key={`${rowIdx}-${colIdx}`}
              type="number"
              step="any"
              value={cell}
              onChange={(e) => editable && onInputChange(rowIdx, colIdx, e.target.value)}
              readOnly={!editable}
              className="matrix-cell"
            />
          ))}
          <span className="matrix-bracket">]</span>
        </div>
      ))}
    </div>
  </div>
));

function App() {
  const [matrix, setMatrix] = useState([
    [0, 0, 0],
    [0, 0, 0],
    [0, 0, 0]
  ]);
  const [matrixInputs, setMatrixInputs] = useState([
    ['', '', ''],
    ['', '', ''],
    ['', '', '']
  ]);
  const [rref, setRref] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleInputChange = useCallback((row, col, value) => {
    // Update the input string state
    setMatrixInputs(prev => {
      return prev.map((r, rIdx) =>
        r.map((c, cIdx) => {
          if (rIdx === row && cIdx === col) {
            return value;
          }
          return c;
        })
      );
    });
    
    // Update the numeric matrix state
    setMatrix(prevMatrix => {
      return prevMatrix.map((r, rIdx) =>
        r.map((c, cIdx) => {
          if (rIdx === row && cIdx === col) {
            const numValue = parseFloat(value);
            return isNaN(numValue) ? 0 : numValue;
          }
          return c;
        })
      );
    });
    
    // Only clear RREF and error if they exist (using functional updates to avoid dependencies)
    setRref(prev => prev !== null ? null : prev);
    setError(prev => prev !== null ? null : prev);
  }, []);

  const handleCalculate = async () => {
    setLoading(true);
    setError(null);
    
    try {
      // Use environment variable for API URL, fallback to proxy for local development
      const apiUrl = process.env.REACT_APP_API_URL || '/api/rref';
      
      if (!apiUrl) {
        throw new Error('API URL is not configured');
      }
      
      const response = await fetch(apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          matrix: {
            data: matrix
          }
        }),
      });

      if (!response.ok) {
        throw new Error(`Server error: ${response.statusText}`);
      }

      const data = await response.json();
      setRref(data.rref.data);
    } catch (err) {
      setError(err.message);
      console.error('Error calculating RREF:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleClear = () => {
    setMatrix([
      [0, 0, 0],
      [0, 0, 0],
      [0, 0, 0]
    ]);
    setMatrixInputs([
      ['', '', ''],
      ['', '', ''],
      ['', '', '']
    ]);
    setRref(null);
    setError(null);
  };


  const RREFMatrixDisplay = ({ matrix, title }) => (
    <div className="matrix-container">
      <h3>{title}</h3>
      <div className="matrix">
        {matrix.map((row, rowIdx) => (
          <div key={rowIdx} className="matrix-row">
            <span className="matrix-bracket">[</span>
            {row.map((cell, colIdx) => (
              <div
                key={colIdx}
                className="matrix-cell-display"
              >
                {cell === 0 ? '0' : cell}
              </div>
            ))}
            <span className="matrix-bracket">]</span>
          </div>
        ))}
      </div>
    </div>
  );

  return (
    <div className="App">
      <div className="container">
        <h1>3×3 Matrix RREF Calculator</h1>
        <p className="subtitle">Enter your matrix values and calculate the Reduced Row Echelon Form</p>
        
        <MatrixDisplay 
          matrixInputs={matrixInputs} 
          onInputChange={handleInputChange} 
          title="Input Matrix" 
          editable={true} 
        />
        
        <div className="button-group">
          <button onClick={handleCalculate} disabled={loading} className="btn btn-primary">
            {loading ? 'Calculating...' : 'Calculate RREF'}
          </button>
          <button onClick={handleClear} className="btn btn-secondary">
            Clear
          </button>
        </div>

        {error && (
          <div className="error-message">
            Error: {error}
          </div>
        )}

        {rref && (
          <RREFMatrixDisplay matrix={rref} title="RREF Result" />
        )}
      </div>
    </div>
  );
}

export default App;

