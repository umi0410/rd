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
import { makeStyles } from '@material-ui/core/styles';
import Box from'@mui/material/Box'
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Grid from'@mui/material/Grid'
import { ThemeProvider } from '@material-ui/core/styles';
import theme from './theme';

const useStyles = makeStyles((theme) => ({
  form: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
  },
  textField: {
    margin: theme.spacing(1),
    width: '100%',
  },
  button: {
    margin: theme.spacing(1),
    width: '100%',
  },
}));

export default function AddAlias() {
  const classes = useStyles();
  
  const [formData, setFormData] = useState({
    group: 'default',
    name: '',
    destination: window.location.toString(),
    description: '',
  });

  const handleChange = (event) => {
    setFormData({
      ...formData,
      [event.target.name]: event.target.value,
    });
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    // Handle form submission logic here
    fetch('http://localhost:18080/aliases', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(formData),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log('Success:', data);
        // handle successful response here
      })
      .catch((error) => {
        console.error('Error:', error);
        // handle error here
      });
  };

  return (
    <form onSubmit={handleSubmit}>
      <ThemeProvider theme={theme}>
    <Grid container>
      <Grid item xs={12}>
        <TextField
          label="Group"
          name="group"
          value={formData.group}
          onChange={handleChange}
          required
          fullWidth
          size="small"
        />
      </Grid>
      <Grid item xs={12}>
        <TextField
          label="Name"
          name="name"
          value={formData.name}
          onChange={handleChange}
          required
          fullWidth
          size="small"
        />
      </Grid>
      <Grid item xs={12}>
        <TextField
          label="Destination"
          name="destination"
          value={formData.destination}
          onChange={handleChange}
          required
          fullWidth
          size="small"
        />
      </Grid>
      <Grid item xs={12}>
        <TextField
          label="Description"
          name="description"
          rows={4}
          value={formData.description}
          onChange={handleChange}
          fullWidth
          size="small"
        />
      </Grid>
      <Grid item xs={12} mt={2}>
        <Button
          variant="contained"
          color="primary"
          type="submit"
          fullWidth
          size="small"
        >
          Submit
        </Button>
      </Grid>
    </Grid>
    </ThemeProvider>
  </form>
  );
}