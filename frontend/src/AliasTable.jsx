import * as React from 'react';
import {useEffect, useState} from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';
import Grid from'@mui/material/Grid'
import Box from'@mui/material/Box'

export default function AliasTable(props) {
  const {isDetailed} = props;
  const [aliases, setAliases] = useState([])

  // Option 1) To use a async function
  useEffect(()=>{
    const numOfSkeletonRows = 32;
    
    // XXX: useEffect is invoked only once
    // and it seems like the same is applied to setAliases().
    // setAliases([...Array(numOfSkeletonRows)].forEach((_, i) => ({
    //   group: "",
    //   name: "",
    //   destination: ""
    // })))
    const fetchData = async () => {
      const result = await fetch("http://localhost:18080/aliases")
      const data = await result.json()

      setAliases(data)
    }
    fetchData()
  }, []);

  // Option 2) To use then method
  // useEffect(()=>{
  //   fetch("http://localhost:18080/aliases")
  //   .then(res => res.json())
  //   .then(data => setAliases(data))
  // }, []);

  return (
    <Box component={Paper}>
      <Table style={{minWidth: '200px'}} aria-label="simple table">
        <TableHead>
          <TableRow >
            {isDetailed &&
              <TableCell align="center"><Typography>Group</Typography></TableCell>
            }

            <TableCell align="center" style={{width: '20%'}}><Typography>Alias</Typography></TableCell>
            <TableCell align="center"><Typography>Destination</Typography></TableCell>
            {isDetailed &&
              <TableCell align="center"><Typography>Description</Typography></TableCell>
            }            
          </TableRow>
        </TableHead>
        <TableBody>

        {aliases.map((alias, i) => (
          <TableRow key={i}>
            {isDetailed &&
              <TableCell sx={{color: 'text.secondary'}} align="center"><Typography>{alias.group}</Typography></TableCell>
            }
            
            <TableCell style={{width: '20%'}}><Typography>{alias.name}</Typography></TableCell>
            <TableCell><Typography><Link href={alias.destination}>{alias.destination}</Link></Typography></TableCell>
            {isDetailed &&
              <TableCell><Typography>{alias.description}</Typography></TableCell>
            }            
          </TableRow>
        ))} 
        </TableBody>
      </Table>
    </Box>
  );
}