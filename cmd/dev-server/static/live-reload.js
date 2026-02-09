<!-- cmd/dev-server/static/live-reload.js -->
(function() {
    // Only in development
    if (window.location.hostname !== 'localhost' && 
        window.location.hostname !== '127.0.0.1') {
        return;
    }
    
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    let reconnectAttempts = 0;
    let maxReconnectAttempts = 10;
    let reconnectDelay = 1000;
    
    function connect() {
        const ws = new WebSocket(wsUrl);
        
        ws.onopen = function() {
            console.log('🔗 Live reload connected');
            reconnectAttempts = 0;
        };
        
        ws.onmessage = function(event) {
            const message = JSON.parse(event.data);
            
            if (message.type === 'reload') {
                console.log(`🔄 Reloading page: ${message.reason}`);
                
                // Try to preserve component state
                const state = window.app?.state;
                if (state) {
                    localStorage.setItem('__dev_state', JSON.stringify(state));
                }
                
                // Reload the page
                setTimeout(() => {
                    window.location.reload();
                }, 100);
            }
        };
        
        ws.onclose = function() {
            console.log('📡 Live reload disconnected');
            
            if (reconnectAttempts < maxReconnectAttempts) {
                reconnectAttempts++;
                const delay = reconnectDelay * Math.pow(1.5, reconnectAttempts);
                
                console.log(`Reconnecting in ${delay}ms...`);
                setTimeout(connect, delay);
            }
        };
        
        ws.onerror = function(error) {
            console.error('Live reload error:', error);
        };
    }
    
    // Restore state on load
    window.addEventListener('load', function() {
        const savedState = localStorage.getItem('__dev_state');
        if (savedState && window.app) {
            try {
                const state = JSON.parse(savedState);
                Object.assign(window.app.state, state);
                localStorage.removeItem('__dev_state');
            } catch (e) {
                console.warn('Failed to restore state:', e);
            }
        }
    });
    
    // Start connection
    connect();
    
    // Export for manual control
    window.liveReload = {
        reconnect: connect,
        disconnect: function() {
            // Implementation for manual disconnection
        }
    };
})();
