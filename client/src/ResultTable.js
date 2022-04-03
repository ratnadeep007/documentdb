import { useEffect, useState } from "react";
import { JsonToTable } from "react-json-to-table";
import AddData from "./AddData";

function ResultTable() {
  let [data, setData] = useState({});

  const makeCall = () => {
    fetch("http://localhost:8080/docs/0")
      .then(response => response.json())
      .then(data => {
        if (data) {
          let values = data['body']['values'];
          let valueList = [];
          for (var i = 0; i < values.length; i++) {
            console.log(values[i]);
            if (values[i]['body']) {
              let keys = Object.keys(values[i]['body']);
              let d = {
                "id": values[i]['id'],
              }

              for (let key of keys) {
                d[key] = values[i]['body'][key];
              }

              valueList.push(d);
            }
          }
          setData(valueList);
        } 
      });
  }

  useEffect(() => {
    makeCall();
    setInterval(() => {
      makeCall();
    }, 5000);
  }, []);

  return (
    <>
      <div className="container mx-auto mt-2">
        <AddData />
        <JsonToTable json={data} />
      </div>
    </>
  )
}

export default ResultTable;