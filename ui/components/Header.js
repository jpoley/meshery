import React from 'react';
import PropTypes from 'prop-types';
import AppBar from '@material-ui/core/AppBar';
import Grid from '@material-ui/core/Grid';
import Hidden from '@material-ui/core/Hidden';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import {connect} from "react-redux";
import User from './User';
import NoSsr from '@material-ui/core/NoSsr';
import MesheryNotification from './MesheryNotification';

const lightColor = 'rgba(255, 255, 255, 0.7)';

const styles = theme => ({
  secondaryBar: {
    zIndex: 0,
  },
  menuButton: {
    marginLeft: -theme.spacing(1),
  },
  iconButtonAvatar: {
    padding: 4,
  },
  link: {
    textDecoration: 'none',
    color: lightColor,
    '&:hover': {
      color: theme.palette.common.white,
    },
  },
  button: {
    borderColor: lightColor,
  },
  notifications: {
    paddingLeft: theme.spacing(4),
    paddingRight: theme.spacing(0),
    marginLeft: theme.spacing(4),
  },
  userContainer: {
    paddingLeft: 1,
    display: 'flex',
  },
  userSpan: {
    marginLeft: theme.spacing(1),
  }
});

class Header extends React.Component {

  render() {
    const { classes, title, onDrawerToggle } = this.props;

    return (
      <NoSsr>
      <React.Fragment>
        <AppBar color="primary" position="sticky" elevation={0}>
          <Toolbar>
            <Grid container spacing={8} alignItems="center">
              <Hidden smUp>
                <Grid item>
                  <IconButton
                    color="inherit"
                    aria-label="Open drawer"
                    onClick={onDrawerToggle}
                    className={classes.menuButton}
                  >
                    <MenuIcon />
                  </IconButton>
                </Grid>
              </Hidden>
              <Grid container xs justify="center">
                <Grid item>
                </Grid>
              </Grid>
              {/* <Grid item className={classes.notifications}>
                <MesheryNotification />
              </Grid> */}
              <Grid item className={classes.userContainer}>
                {/* <IconButton color="inherit" className={classes.iconButtonAvatar}>
                  <Avatar className={classes.avatar} src="/static/images/avatar/1.jpg" />
                </IconButton> */}
                <MesheryNotification />
                <span className={classes.userSpan}>
                <User color="inherit" iconButtonClassName={classes.iconButtonAvatar} avatarClassName={classes.avatar} />
                </span>
              </Grid>
            </Grid>
          </Toolbar>
        </AppBar>
        <AppBar
          component="div"
          className={classes.secondaryBar}
          color="primary"
          position="static"
          elevation={0}
        >
          <Toolbar>
            <Grid container alignItems="center" spacing={8}>
              <Grid item xs>
                <Typography color="inherit" variant="h5">
                  {title}
                </Typography>
              </Grid>
              {/* <Grid item>
                <Button className={classes.button} variant="outlined" color="inherit" size="small">
                  Web setup
                </Button>
              </Grid> */}
              {/* <Grid item>
                <Tooltip title="Help">
                  <IconButton color="inherit">
                    <HelpIcon />
                  </IconButton>
                </Tooltip>
              </Grid> */}
            </Grid>
          </Toolbar>
        </AppBar>
        {/* <AppBar
          component="div"
          className={classes.secondaryBar}
          color="primary"
          position="static"
          elevation={0}
        >
          <Tabs value={0} textColor="inherit">
            <Tab textColor="inherit" label="Users" />
            <Tab textColor="inherit" label="Sign-in method" />
            <Tab textColor="inherit" label="Templates" />
            <Tab textColor="inherit" label="Usage" />
          </Tabs>
        </AppBar> */}
        
      </React.Fragment>
      </NoSsr>
    );
  }
}

Header.propTypes = {
  classes: PropTypes.object.isRequired,
  onDrawerToggle: PropTypes.func.isRequired,
};

const mapStateToProps = state => {
  // console.log("header - mapping state to props. . . new title: "+ state.get("page").get("title"));
  // console.log("state: " + JSON.stringify(state));
  return { title: state.get("page").get("title") }
}

// const mapDispatchToProps = dispatch => {
//   return {
//     updatePageAndTitle: bindActionCreators(updatePageAndTitle, dispatch)
//   }
// }

export default withStyles(styles)(connect(
  mapStateToProps
)(Header));