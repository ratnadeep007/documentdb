import { useEffect, useState } from 'react';

function Status() {
  let [serverOnline, setServerOnline] = useState(false);

  useEffect(() => {
    fetch("http://localhost:8080/status")
      .then(response => response.status)
      .then(data => setServerOnline(data))
  })

  return (
    <div className="container mx-auto mt-2 flex-row">
      <div className="text-lg font-semibold">Status:
        { 
          serverOnline === 200 ? 
            <span className="text-green-600"> Online</span> : 
            <span className="text-red-600"> Offline</span>
        }
      </div>
    </div>
  )
}

export default Status;