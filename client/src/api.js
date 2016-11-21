const WSS_ADDR = process.env.API_ADDR || `ws://${window.location.host}`

export const newSession = function(sessionId){
  return new Promise((resolve, reject) => {
    var socket = new WebSocket(`${WSS_ADDR}/session/${sessionId}`);
    socket.onopen = () => resolve(socket)
    socket.onerror = (e) => reject(e)
  });
}
