import * as React from 'react';
import PropTypes from 'prop-types';
import { alpha } from '@mui/material/styles';
import Box from '@mui/material/Box';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TablePagination from '@mui/material/TablePagination';
import TableRow from '@mui/material/TableRow';
import TableSortLabel from '@mui/material/TableSortLabel';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import Checkbox from '@mui/material/Checkbox';
import IconButton from '@mui/material/IconButton';
import Tooltip from '@mui/material/Tooltip';
import FormControlLabel from '@mui/material/FormControlLabel';
import Switch from '@mui/material/Switch';
import DeleteIcon from '@mui/icons-material/Delete';
import FilterListIcon from '@mui/icons-material/FilterList';
import { visuallyHidden } from '@mui/utils';

function createData(sourcePod, desPod, desIP, sourcePort, desPort, namespace, status, lastSeen) 
{
  return {sourcePod, desPod, desIP, sourcePort, desPort, namespace, status, lastSeen};
}



const jsonData = [{"IP":{"destination":"10.244.1.207","ipVersion":"IPv4","source":"192.168.56.1"},"Summary":"TCP Flags: ACK","Type":"L3_L4","destination":{"ID":900,"identity":11848,"labels":["k8s:app.kubernetes.io/component=controller","k8s:app.kubernetes.io/instance=ingress-nginx","k8s:app.kubernetes.io/name=ingress-nginx","k8s:io.cilium.k8s.namespace.labels.app.kubernetes.io/instance=ingress-nginx","k8s:io.cilium.k8s.namespace.labels.app.kubernetes.io/name=ingress-nginx","k8s:io.cilium.k8s.namespace.labels.kubernetes.io/metadata.name=ingress-nginx","k8s:io.cilium.k8s.policy.cluster=default","k8s:io.cilium.k8s.policy.serviceaccount=ingress-nginx","k8s:io.kubernetes.pod.namespace=ingress-nginx"],
"namespace":"ingress-nginx","pod_name":"ingress-nginx-controller-6c965d68c9-wvk8n","workloads":[{"kind":"ReplicaSet","name":"ingress-nginx-controller-6c965d68c9"}]},

"ethernet":{"destination":"9a:11:9d:ef:21:b4","source":"f2:97:91:2d:29:3a"},"event_type":{"type":4},"interface":{"index":8,"name":"lxcf88934bbe643"},"is_reply":false, "xxx":{"TCP":{"destination_port":80,"flags":{"ACK":true},"source_port":44354}},"node_name":"kind-worker","source":{"identity":2,"labels":["reserved:world"]},"time":"2022-06-27T12:17:10.105255515Z","trace_observation_point":"TO_ENDPOINT","traffic_direction":"INGRESS","verdict":"FORWARDED"}];

const getRow = (data) => {
  let rows = [];
  data.forEach(flow => {
    let row = createData(flow.IP.source, flow.destination.pod_name, flow.IP.destination, flow["xxx"].TCP.source_port, flow["xxx"].TCP.destination_port, flow.destination.namespace, flow.verdict, flow.time);
    rows.push(row);
  });

  return rows;
}

const rows = getRow(jsonData);

function descendingComparator(a, b, orderBy) {
  if (b[orderBy] < a[orderBy]) {
    return -1;
  }
  if (b[orderBy] > a[orderBy]) {
    return 1;
  }
  return 0;
}

function getComparator(order, orderBy) {
  return order === 'desc'
    ? (a, b) => descendingComparator(a, b, orderBy)
    : (a, b) => -descendingComparator(a, b, orderBy);
}

// This method is created for cross-browser compatibility, if you don't
// need to support IE11, you can use Array.prototype.sort() directly
function stableSort(array, comparator) {
  const stabilizedThis = array.map((el, index) => [el, index]);
  stabilizedThis.sort((a, b) => {
    const order = comparator(a[0], b[0]);
    if (order !== 0) {
      return order;
    }
    return a[1] - b[1];
  });
  return stabilizedThis.map((el) => el[0]);
}

const headCells = [
  {
    id: 'sourcePod',
    numeric: false,
    disablePadding: true,
    label: 'Source Pod Name/Source IP',
  },
  {
    id: 'desPod',
    numeric: false,
    disablePadding: false,
    label: 'Destination Pod Name',
  },
  {
    id: 'desIP',
    numeric: false,
    disablePadding: false,
    label: 'Destionation IP',
  },
  {
    id: 'sourcePort',
    numeric: false,
    disablePadding: false,
    label: 'Source Port',
  },
  {
    id: 'desPort',
    numeric: false,
    disablePadding: false,
    label: 'Destination Port',
  },
  {
    id: 'namespace',
    numeric: false,
    disablePadding: true,
    label: 'Namespace',
  },
  {
    id: 'status',
    numeric: false,
    disablePadding: false,
    label: 'Status',
  },
  {
    id: 'last-seen',
    numeric: false,
    disablePadding: false,
    label: 'Last Seen',
  },
];

function EnhancedTableHead(props) {
  const { onSelectAllClick, order, orderBy, numSelected, rowCount, onRequestSort } =
    props;
  const createSortHandler = (property) => (event) => {
    onRequestSort(event, property);
  };

  return (
    <TableHead>
      <TableRow>
        <TableCell padding="checkbox"
            sx={{
                backgroundColor: "aqua",
                borderBottom: "2px solid black",
                borderTop: "2px solid black",
                "& th": {
                  fontSize: "1.25rem",
                  color: "rgba(96, 96, 96)"
                }
              }}
            >
            
          <Checkbox
            color="primary"
            indeterminate={numSelected > 0 && numSelected < rowCount}
            checked={rowCount > 0 && numSelected === rowCount}
            onChange={onSelectAllClick}
            inputProps={{
              'aria-label': 'select all desserts',
            }}
          />
        </TableCell>
        {headCells.map((headCell) => (
          <TableCell
            key={headCell.id}
            align={'left'}
            sortDirection={orderBy === headCell.id ? order : false}
            sx={{
                paddingLeft: "4px",
                paddingRight: "0px",
                paddingY: "2px",
                backgroundColor: "aqua",
                borderBottom: "2px solid black",
                borderTop: "2px solid black",
                borderRight: "1px solid black",
                "& th": {
                  fontSize: "1.25rem",
                  color: "rgba(96, 96, 96)"
                }
              }}
          >
            <TableSortLabel
              active={orderBy === headCell.id}
              direction={orderBy === headCell.id ? order : 'asc'}
              onClick={createSortHandler(headCell.id)}
            >
              {headCell.label}
              {orderBy === headCell.id ? (
                <Box component="span" sx={visuallyHidden}>
                  {order === 'desc' ? 'sorted descending' : 'sorted ascending'}
                </Box>
              ) : null}
            </TableSortLabel>
          </TableCell>
        ))}
      </TableRow>
    </TableHead>
  );
}

EnhancedTableHead.propTypes = {
  numSelected: PropTypes.number.isRequired,
  onRequestSort: PropTypes.func.isRequired,
  onSelectAllClick: PropTypes.func.isRequired,
  order: PropTypes.oneOf(['asc', 'desc']).isRequired,
  orderBy: PropTypes.string.isRequired,
  rowCount: PropTypes.number.isRequired,
};

const EnhancedTableToolbar = (props) => {
  const { numSelected } = props;

  return (
    <Toolbar
      sx={{
        pl: { sm: 2 },
        pr: { xs: 1, sm: 1 },
        ...(numSelected > 0 && {
          bgcolor: (theme) =>
            alpha(theme.palette.primary.main, theme.palette.action.activatedOpacity),
        }),
      }}
    >
      {numSelected > 0 ? (
        <Typography
          sx={{ flex: '1 1 100%' }}
          color="inherit"
          variant="subtitle1"
          component="div"
        >
          {numSelected} selected
        </Typography>
      ) : (
        <Typography
          sx={{ flex: '1 1 100%' }}
          variant="h6"
          id="tableTitle"
          component="div"
        >
          Suspected Flow
        </Typography>
      )}

      {numSelected > 0 ? (
        <Tooltip title="Delete">
          <IconButton>
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      ) : (
        <Tooltip title="Filter list">
          <IconButton>
            <FilterListIcon />
          </IconButton>
        </Tooltip>
      )}
    </Toolbar>
  );
};

EnhancedTableToolbar.propTypes = {
  numSelected: PropTypes.number.isRequired,
};

export default function EnhancedTable() {
  const [order, setOrder] = React.useState('asc');
  const [orderBy, setOrderBy] = React.useState('calories');
  const [selected, setSelected] = React.useState([]);
  const [page, setPage] = React.useState(0);
  const [dense, setDense] = React.useState(false);
  const [rowsPerPage, setRowsPerPage] = React.useState(5);

  const handleRequestSort = (event, property) => {
    const isAsc = orderBy === property && order === 'asc';
    setOrder(isAsc ? 'desc' : 'asc');
    setOrderBy(property);
  };

  const handleSelectAllClick = (event) => {
    if (event.target.checked) {
      const newSelecteds = rows.map((n) => n.name);
      setSelected(newSelecteds);
      return;
    }
    setSelected([]);
  };

  const handleClick = (event, name) => {
    const selectedIndex = selected.indexOf(name);
    let newSelected = [];

    if (selectedIndex === -1) {
      newSelected = newSelected.concat(selected, name);
    } else if (selectedIndex === 0) {
      newSelected = newSelected.concat(selected.slice(1));
    } else if (selectedIndex === selected.length - 1) {
      newSelected = newSelected.concat(selected.slice(0, -1));
    } else if (selectedIndex > 0) {
      newSelected = newSelected.concat(
        selected.slice(0, selectedIndex),
        selected.slice(selectedIndex + 1),
      );
    }

    setSelected(newSelected);
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleChangeDense = (event) => {
    setDense(event.target.checked);
  };

  const isSelected = (name) => selected.indexOf(name) !== -1;

  // Avoid a layout jump when reaching the last page with empty rows.
  const emptyRows =
    page > 0 ? Math.max(0, (1 + page) * rowsPerPage - rows.length) : 0;

  return (
    <Box sx={{ width: '100%' }}>
      <Paper sx={{ width: '100%', mb: 2 }}>
        <EnhancedTableToolbar numSelected={selected.length} />
        <TableContainer>
          <Table
            sx={{ minWidth: 750 }}
            aria-labelledby="tableTitle"
            size={dense ? 'small' : 'medium'}
          >
            <EnhancedTableHead
              numSelected={selected.length}
              order={order}
              orderBy={orderBy}
              onSelectAllClick={handleSelectAllClick}
              onRequestSort={handleRequestSort}
              rowCount={rows.length}
            />
            <TableBody>
              {/* if you don't need to support IE11, you can replace the `stableSort` call with:
                 rows.slice().sort(getComparator(order, orderBy)) */}
              {stableSort(rows, getComparator(order, orderBy))
                .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                .map((row, index) => {
                  const isItemSelected = isSelected(row.name);
                  const labelId = `enhanced-table-checkbox-${index}`;

                  return (
                    <TableRow
                      hover
                      onClick={(event) => handleClick(event, row.name)}
                      role="checkbox"
                      aria-checked={isItemSelected}
                      tabIndex={-1}
                      key={row.name}
                      size={'small'}
                      selected={isItemSelected}
                    >
                      <TableCell padding="checkbox">
                        <Checkbox
                          color="primary"
                          checked={isItemSelected}
                          inputProps={{
                            'aria-labelledby': labelId,
                          }}
                        />
                      </TableCell>
                      <TableCell
                        component="th"
                        id={labelId}
                        scope="row"
                        padding="none"
                        sx={{
                            padding: "0px 0px",
                            borderRight: "1px solid black",
                          }}
                      >
                        {row.sourcePod}
                      </TableCell>

                      <TableCell
                        align="left" 
                        sx={{
                            padding: "0px 10px",
                            borderRight: "1px solid black",
                          }}
                        >
                          {row.desPod}
                      </TableCell>
                    
                      <TableCell
                        align="left" 
                        sx={{
                            padding: "0px 10px",
                            borderRight: "1px solid black",
                          }}
                        >
                          {row.desIP}
                      </TableCell>

                      <TableCell
                        align="left" 
                        sx={{
                            padding: "0px 10px",
                            borderRight: "1px solid black",
                          }}
                        >
                          {row.sourcePort}
                      </TableCell>

                      <TableCell
                        align="left" 
                        sx={{
                            padding: "0px 10px",
                            borderRight: "1px solid black",
                          }}
                        >
                          {row.desPort}
                      </TableCell>

                      <TableCell
                        align="left" 
                        sx={{
                            padding: "0px 10px",
                            borderRight: "1px solid black",
                          }}
                        >
                          {row.namespace}
                      </TableCell>

                      <TableCell
                        align="left"

                        sx={{
                            padding: "0px 10px",
                            borderRight: "1px solid black",
                            color: 'green'
                          }}
                        >
                          {row.status}
                      </TableCell>

                      <TableCell
                        align="left" 
                        sx={{
                            padding: "0px 10px",
                            borderRight: "1px solid black",
                          }}
                        >
                          {row.lastSeen}
                      </TableCell>
                      
                    </TableRow>
                  );
                })}
              {emptyRows > 0 && (
                <TableRow
                  style={{
                    height: (dense ? 33 : 53) * emptyRows,
                  }}
                >
                  <TableCell colSpan={6} />
                </TableRow>
              )}
            </TableBody>
          </Table>
        </TableContainer>
        <TablePagination
          rowsPerPageOptions={[5, 10, 25]}
          component="div"
          count={rows.length}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </Paper>
      <FormControlLabel
        control={<Switch checked={dense} onChange={handleChangeDense} />}
        label="Dense padding"
      />
    </Box>
  );
}
