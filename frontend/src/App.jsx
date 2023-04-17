import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import Typography from '@mui/material/Typography';
import AliasTable from './AliasTable'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="App">
      <Typography variant="h3" gutterBottom={true}>rd: redirecting by using aliases</Typography>
      <AliasTable group="default"></AliasTable>
    </div>
  )
}

export default App
