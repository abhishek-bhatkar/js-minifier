/*!
 * Complex JavaScript Test File
 * Copyright (c) 2024
 * Licensed under MIT
 */

// Utility function for calculation
function calculateComplex(x, y, z) {
    const multiplier = 2.5;  // Default multiplier
    let result = 0;

    /* This is a multi-line comment
       that should be removed during minification
       but contains important context */

    // Calculate the result
    if (x > y) {
        result = (x + y) * multiplier;
    } else {
        result = (y - x) * multiplier + z;
    }

    return result;
}

// Class definition for data processing
class DataProcessor {
    constructor(initialValue) {
        this.value = initialValue;
        this.history = [];
    }

    // Process the data with given parameters
    processData(inputData) {
        const processedValue = calculateComplex(
            this.value,
            inputData,
            10
        );
        
        this.history.push({
            input: inputData,
            output: processedValue,
            timestamp: new Date()
        });

        return processedValue;
    }

    // Get processing history
    getHistory() {
        return this.history.map(item => {
            return {
                input: item.input,
                output: item.output,
                time: item.timestamp.toISOString()
            };
        });
    }
}
