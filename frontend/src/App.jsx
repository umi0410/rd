import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Typography from '@mui/material/Typography';
import AliasTable from './AliasTable'
import AddAlias from './AddAlias'
import Box from'@mui/material/Box'
import Grid from'@mui/material/Grid'

function App() {
  const [count, setCount] = useState(0)

  return (
      <Box className="root" padding={0}>
        {/* TODO: Can I use Reac Router? because popup file is set to "index.html" */}
        {/* <BrowserRouter>
          <Routes>
            <Route path="/popup"></Route>
            <Route path="/*" element={<Grid container><Typography variant="h3" gutterBottom={true}>rd: redirecting by using aliases</Typography>
              </Grid>}></Route>
          </Routes>

        <Grid container>
          <AliasTable></AliasTable>
        </Grid>  
        </BrowserRouter> */}
        <Box>
          <AddAlias></AddAlias>
        </Box>
        <Box>
          <AliasTable props={{isDetailed: false}}></AliasTable>
        </Box>
      </Box>
  )
}

export default App
