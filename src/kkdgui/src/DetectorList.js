import * as React from 'react';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ListSubheader from '@mui/material/ListSubheader';
import Switch from '@mui/material/Switch';
import WifiIcon from '@mui/icons-material/Wifi';
import BluetoothIcon from '@mui/icons-material/Bluetooth';
import SettingsInputComponentIcon from '@mui/icons-material/SettingsInputComponent';
import { Typography } from '@mui/material';
import AddCircleIcon from '@mui/icons-material/AddCircle';

export default function DetectorList() {
  const [checked, setChecked] = React.useState(['wifi']);
  const [names, setNames] = React.useState([]);

  const handleToggle = (value) => () => {
    const currentIndex = checked.indexOf(value);
    const newChecked = [...checked];

    if (currentIndex === -1) {
      newChecked.push(value);
    } else {
      newChecked.splice(currentIndex, 1);
    }

    setChecked(newChecked);
  };

  React.useEffect(async () => {

    const res = await fetch('/apinode/getDetectors');
    const names = await res.json();

    const nameList = names.nameList;
    setNames(nameList);
  }, []);

  return (

    <List
      sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}
    >
      {names.map(name => {

        return (
          <ListItem>
            <ListItemIcon>
              <SettingsInputComponentIcon />
            </ListItemIcon>
            <ListItemText 
              primary={<Typography variant="h6" style={{ color: 'green' }}>{name}</Typography>}
              secondary="Added: 2 days ago"
            />

            <Switch
              edge="end"
              onChange={handleToggle('wifi')}
              checked={checked.indexOf('wifi') !== -1}
              sx={{marginLeft: 3}}
              inputProps={{
                'aria-labelledby': 'switch-list-label-wifi',
              }}
            />
          </ListItem>
        )
      })}
    </List>
  );
}
