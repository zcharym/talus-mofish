import { useState } from 'react'
import {GreetService, AppService} from "../bindings/github.com/songwei.ma/talus_echo_loop";

function App() {
  const [name, setName] = useState<string>('');
  const [result, setResult] = useState<string>('Talus Echo — English learning (infrastructure ready)');
  const [dbPath, setDbPath] = useState<string>('');

  const doGreet = () => {
    const localName = name || 'anonymous';
    GreetService.Greet(localName).then((resultValue: string) => {
      setResult(resultValue);
    }).catch((err: unknown) => {
      console.error(err);
    });
  }

  const loadDbPath = () => {
    AppService.DatabasePath().then((path: string) => {
      setDbPath(path);
    }).catch((err: unknown) => {
      console.error(err);
    });
  }

  return (
    <div className="container">
      <div>
        <a data-wml-openURL="https://wails.io">
          <img src="/wails.png" className="logo" alt="Wails logo"/>
        </a>
        <a data-wml-openURL="https://reactjs.org">
          <img src="/react.svg" className="logo react" alt="React logo"/>
        </a>
      </div>
      <h1>Talus Echo</h1>
      <div className="result">{result}</div>
      <div className="card">
        <div className="input-box">
          <input className="input" value={name} onChange={(e) => setName(e.target.value)} type="text" autoComplete="off" placeholder="Your name"/>
          <button className="btn" onClick={doGreet}>Greet</button>
        </div>
        <div className="input-box" style={{ marginTop: '1rem' }}>
          <button className="btn" type="button" onClick={loadDbPath}>Show database path</button>
        </div>
        {dbPath ? <p className="hint">SQLite: {dbPath}</p> : null}
      </div>
    </div>
  )
}

export default App
