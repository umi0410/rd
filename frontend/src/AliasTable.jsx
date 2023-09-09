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
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import DeleteIcon from '@mui/icons-material/Delete';
import Grid from'@mui/material/Grid'
import Box from'@mui/material/Box'

export default function AliasTable(props) {
  const {isDetailed} = props;
  const [aliases, setAliases] = useState([])

  const fetchData = async () => {
      const result = await fetch("http://localhost:18080/aliases")
      const data = await result.json()

      setAliases(data)
    }

  useEffect(()=>{
    const numOfSkeletonRows = 32;
    fetchData()
  }, []);

  const deleteAlias = async (id) =>{
    const result = await fetch("http://localhost:18080/aliases",
        {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id: id
        })
    })
    await fetchData()
  }
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
            <TableCell align="center"><Typography>Actions</Typography></TableCell>
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
            <TableCell><IconButton aria-label="delete"  variant="outlined" onClick={e=>{deleteAlias(alias.id)}}> <DeleteIcon fontSize="small" /> </IconButton></TableCell>
          </TableRow>
        ))} 
        </TableBody>
      </Table>
    </Box>
  );
}
