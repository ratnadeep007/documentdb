import Editor from "@monaco-editor/react";
import { useState } from "react";

function AddData() {
  let [data, setData] = useState();

  const submitData = (data) => {
    fetch("http://localhost:8080/docs", {
      method: "POST",
      mode: "cors",
      headers: {
        'Content-Type': 'application/json',
      },
      body: data
    }).then(response => response.json())
      .then(response => {
        if (response["body"]["id"]) {
          setData("")
        }
      })
  }

  return (
    <div className="mt-2 flex-col mb-2">
      <div className="text-gray-700">Enter your JSON here to add</div>
      <Editor
        height="30vh"
        width="90vw"
        defaultLanguage="json"
        defaultValue=""
        onChange={(value) => setData(value)}
        value={data}
      />
      <div 
        className="bg-green-600 rounded-md w-12 flex justify-center text-white hover:cursor-pointer"
        onClick={() => submitData(data)}
      >Add</div>
    </div>
  )
}

export default AddData;