import * as React from 'react';
import PropTypes from 'prop-types';
import SwipeableViews from 'react-swipeable-views';
import { useTheme } from '@mui/material/styles';
import AppBar from '@mui/material/AppBar';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import { Button, Grid } from '@mui/material';
import EnhancedTable from './Table.js';
import PolicyListAccordion from './PolicesList.js';
import DetectorList from './DetectorList.js';
import { Container } from '@mui/system';
import AddCircle from '@mui/icons-material/AddCircle';
import FlowGraph from './FlowGraph.js';
import TextField from '@mui/material/TextField';


function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`full-width-tabpanel-${index}`}
      aria-labelledby={`full-width-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
}

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.number.isRequired,
  value: PropTypes.number.isRequired,
};

function a11yProps(index) {
  return {
    id: `full-width-tab-${index}`,
    'aria-controls': `full-width-tabpanel-${index}`,
  };
}

const addDetectorHandler = (name, imageName) => {
  fetch(`/apinode/addDetector?name=${name}&imageName=${imageName}`)
  .then(res => {
    alert("Add Detector successfully!");
    //window.location.reload();
  })
}

export default function FullWidthTabs() {
  const theme = useTheme();
  const [value, setValue] = React.useState(0);
  const [name, setName] = React.useState(0);
  const [imageName, setImageName] = React.useState(0);

  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  const handleChangeIndex = (index) => {
    setValue(index);
  };

  return (
    <Box sx={{ bgcolor: 'background.paper', width: 'auto' }}>
      <AppBar position="static">
        <Tabs
          value={value}
          onChange={handleChange}
          indicatorColor="secondary"
          textColor="inherit"
          variant="fullWidth"
          aria-label="full width tabs example"
        >
          <Tab label="Flow and Policies" {...a11yProps(0)} />
          <Tab label="Detectors" {...a11yProps(1)} />
        </Tabs>
      </AppBar>
      <SwipeableViews
        axis={theme.direction === 'rtl' ? 'x-reverse' : 'x'}
        index={value}
        onChangeIndex={handleChangeIndex}
        se
      >
        <TabPanel value={value} index={0} dir={theme.direction}>
          <Grid container spacing={2}>
            <Grid item xs={8}>
              <FlowGraph/>
            </Grid>
            <Grid item xs={4}>
              <PolicyListAccordion></PolicyListAccordion>
            </Grid>

          </Grid>
        </TabPanel>

        <TabPanel value={value} index={1} dir={theme.direction} sx={{m: 60}}>
          <Grid
            container
            spacing={0}
            direction="column"
            alignItems="center"
            justifyContent="center"
            textAlign="center"
          >

            <Grid item xs={3}>
              <DetectorList />
              
              <Box
                display="flex"
                alignItems="center"
                marginTop={2}
              >
              
              
                <TextField id="outlined-basic" label="Detector Name" variant="outlined" 
                    onChange={(e) => setName(e.target.value)}
                    />
                <TextField id="outlined-basic" label="Image URL" variant="outlined"
                    onChange={(e) => setImageName(e.target.value)}
                    />

                <Button variant="contained" 
                  color="success"
                  size="large" startIcon={<AddCircle />}
                  onClick={() => addDetectorHandler(name, imageName)}
                  >
                  Add New Detector
                </Button>

              </Box>

            </Grid>   
          
          </Grid> 
        </TabPanel>
      </SwipeableViews>
    </Box>
  );
}
