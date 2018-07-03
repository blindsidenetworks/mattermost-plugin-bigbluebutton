const React = window.react;
import PropTypes from 'prop-types';
import {makeStyleFromTheme,changeOpacity} from 'mattermost-redux/utils/theme_utils';
import { Link } from 'react-router-dom'
import * as ChannelActions from 'mattermost-redux/actions/channels';


// const dispatch = window.store.dispatch;
// const getState = window.store.getState;
//test test 3
export default class ProfilePopover extends React.PureComponent {
    static propTypes = {

        /*
         * Source URL from the image to display in the popover
         */
        src: PropTypes.string.isRequired,


        /*
         * User the popover is being opened for
         */
        user: PropTypes.object.isRequired,
        state: PropTypes.object.isRequired,
        cur_user:  PropTypes.object.isRequired,
        //directChannel: PropTypes.object.isRequired,

        teamname: PropTypes.string.isRequired,
        /*
         * Status for the user, either 'offline', 'away' or 'online'
         */
        status: PropTypes.string,

        /*
         * Set to true if the user is in a WebRTC call
         */
        isBusy: PropTypes.bool,

        /*
         * Function to call to hide the popover
         */
        hide: PropTypes.func,

        /*
         * Set to true if the popover was opened from the right-hand
         * sidebar (comment thread, search results, etc.)
         */
        isRHS: PropTypes.bool,

        /*
         * Logged in user's theme
         */
        theme: PropTypes.object.isRequired,

        /*
         * The CSS absolute left position
         */
        positionLeft: PropTypes.number.isRequired,

        /*
         * The CSS absolute top position
         */
        positionTop: PropTypes.number.isRequired,

        /* Add custom props here */

        /* Define action props here or remove if no actions */
        actions: PropTypes.shape({

          startMeeting: PropTypes.func.isRequired

        }).isRequired

    }

    static defaultProps = {
        isBusy: false,
        hide: () => {},
        isRHS: false
        /* If necessary, add defaults for custom props here */
    }

    constructor(props) {
        super(props);
      //  this.handleShowDirectChannel = this.handleShowDirectChannel.bind(this);
        //this.handleDirectMessage = this.handleDirectMessage.bind(this);
        this.state = {
        //  currentUserId: UserStore.getCurrentId(),
          loadingDMChannel: -1,
        };
    }

    exampleFunction = () => {
        // Do some things
    }


    handleDirectMessage = async () => {
      const dispatch = window.store.dispatch;
      const result = await ChannelActions.createDirectChannel(this.props.user.id,this.props.cur_user.id)(dispatch, this.props.state);
      console.log("some result:" + result.data.id)

      await this.props.actions.startMeeting(result.data.id, "",this.props.cur_user.username + " "+ this.props.user.username );

    }

    /* Construct and return the JSX to render here. Make sure that rendering is solely based
        on props and state. */
    render() {
        const style = getStyle(this.props.theme);
        const user = this.props.user;

        const myteam = this.props.teamname
        const url = '/' + myteam + '/messages/@' + user.username
        //const team = this.props.teamName

        return (
          <div
              style={{...style.container, left: this.props.positionLeft, top: this.props.positionTop}}
          >
              <h3 style={style.title}><a>{user.username}</a></h3>
              <div style={style.content}>
                  <img
                      style={style.img}
                      src={this.props.src}
                  />
                  <div style={style.fullName}>
                      {user.first_name + ' ' + user.last_name}
                  </div>
                  <hr style={{margin: '10px -15px 10px'}}/>
                  {this.props.user.id != this.props.cur_user.id &&
                    <div>
                  <Link to= {url}
                    onClick={this.handleDirectMessage}>
                        <i className='fa fa-video-camera'/>{'  Start BigBlueButton Meeting'}
                    </Link>  <br />
                    </div>
                }
                <Link to= {url}>
                      <i className='fa fa-paper-plane'/>{' Send Message'}
                  </Link>
              </div>
          </div>
        );
    }
}

/* Define CSS styles here */
const getStyle = makeStyleFromTheme((theme) => {
    return {
      container: {
          backgroundColor: theme.centerChannelBg,
          position: 'absolute',
          border: '1px solid ' + changeOpacity(theme.centerChannelColor, 0.2),
          borderRadius: '4px',
          zIndex: 9999 // Bring popover to top
      },
      title: {
          padding: '8px 14px',
          margin: '0',
          fontSize: '14px',
          backgroundColor: changeOpacity(theme.centerChannelBg, 0.2),
          borderBottom: '1px solid #ebebeb',
          borderRadius: '5px 5px 0 0'
      },
      content: {
          padding: '9px 14px'
      },
      img: {
          verticalAlign: 'middle',
          maxWidth: '100%',
          borderRadius: '128px',
          margin: '0 0 10px'
      },
      fullName: {
          overflow: 'hidden',
          paddingBottom: '7px',
          whiteSpace: 'nowrap',
          textOverflow: 'ellipsis'
      }
    };
});
