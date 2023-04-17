import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import AliasTable from './AliasTable'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="App">
      <AliasTable group="default"></AliasTable>
    </div>
  )
}

export default App
