// Closure and scope test

function createCounter(initialValue = 0) {
    let count = initialValue;
    
    function increment(step = 1) {
        count += step;
        return count;
    }
    
    function decrement(step = 1) {
        count -= step;
        return count;
    }
    
    function reset() {
        count = initialValue;
        return count;
    }
    
    return {
        increment,
        decrement,
        reset,
        get current() {
            return count;
        }
    };
}

// Factory function with private state
function createLogger(prefix) {
    const logHistory = [];
    const maxHistory = 100;
    
    function formatMessage(level, message) {
        const timestamp = new Date().toISOString();
        return `[${timestamp}] ${prefix} - ${level}: ${message}`;
    }
    
    function addToHistory(entry) {
        logHistory.push(entry);
        if (logHistory.length > maxHistory) {
            logHistory.shift();
        }
    }
    
    return {
        log(message) {
            const entry = formatMessage('INFO', message);
            console.log(entry);
            addToHistory(entry);
        },
        
        error(message) {
            const entry = formatMessage('ERROR', message);
            console.error(entry);
            addToHistory(entry);
        },
        
        getHistory() {
            return [...logHistory];
        }
    };
}
