//source: https://stackoverflow.com/questions/56588625/react-show-material-ui-tooltip-only-for-text-that-has-ellipsis/56589224#56589224

import { Tooltip, Typography } from '@mui/material';
import { Box } from '@mui/system';
import React, { useRef, useEffect, useState } from 'react';

const OverflowTip = ({value, mainColor, hoverColor}) => {
  // Create Ref
  const textElementRef = useRef();

  const compareSize = () => {
    const compare =
      textElementRef.current.scrollWidth > textElementRef.current.clientWidth;
    console.log('compare: ', compare);
    setHover(compare);
  };

  // compare once and add resize listener on "componentDidMount"
  useEffect(() => {
    compareSize();
    window.addEventListener('resize', compareSize);
  }, []);

  // remove resize listener again on "componentWillUnmount"
  useEffect(() => () => {
    window.removeEventListener('resize', compareSize);
  }, []);

  // Define state and function to update the value
  const [hoverStatus, setHover] = useState(false);

  return (
    <Tooltip
      title={value}
      interactive
      disableHoverListener={!hoverStatus}
      style={{fontSize: '2em', color: 'white'}}
    >
      <Box
        ref={textElementRef}
        sx={{color: 'white',
          borderRadius: '16px',
        }}
        margin={0.4}
        padding={0.2}
        paddingX={0.8}
        style={{
          whiteSpace: 'nowrap',
          overflow: 'hidden',
          textOverflow: 'ellipsis',          
          backgroundColor: mainColor,
          '&:hover': {
            backgroundColor: hoverColor,
            opacity: [0.9, 0.8, 0.7],
        },
        }}
      >
        <Typography variant="caption" color="white">
          {value}
        </Typography>
        
      </Box>

    </Tooltip>
  );
};

export default OverflowTip;