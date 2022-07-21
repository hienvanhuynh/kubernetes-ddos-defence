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

  return (
    <List
      sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}
    >
      <ListItem>
        <ListItemIcon>
          <SettingsInputComponentIcon />
        </ListItemIcon>
        <ListItemText 
          primary={<Typography variant="h6" style={{ color: 'green' }}>Detector 1</Typography>}
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
      

      <ListItem>
        <ListItemIcon>
          <SettingsInputComponentIcon />
        </ListItemIcon>
        <ListItemText 
          primary={<Typography variant="h6" style={{ color: 'green' }}>Detector 2</Typography>}
          secondary="Added: 2 days ago"
        />

        <Switch
          edge="end"
          onChange={handleToggle('wifi2')}
          checked={checked.indexOf('wifi2') !== -1}
          sx={{marginLeft: 3}}
          inputProps={{
            'aria-labelledby': 'switch-list-label-wifi',
          }}
        />
      </ListItem>

      <ListItem>
        <ListItemIcon>
          <SettingsInputComponentIcon />
        </ListItemIcon>
        <ListItemText 
          primary={<Typography variant="h6" style={{ color: 'grey' }}>Detector 3</Typography>}
          secondary="Added: 2 days ago"
        />

        <Switch
          edge="end"
          onChange={handleToggle('wifi3')}
          checked={checked.indexOf('wifi3') !== -1}
          sx={{marginLeft: 3}}
          inputProps={{
            'aria-labelledby': 'switch-list-label-wifi',
          }}
        />
      </ListItem>
    </List>
  );
}
