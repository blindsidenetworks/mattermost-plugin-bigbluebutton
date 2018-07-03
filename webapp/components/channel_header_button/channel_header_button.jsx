const React = window.react;
const {Overlay, OverlayTrigger, Popover, Tooltip} = window['react-bootstrap'];


import PopoverListMembersItem from './popover_list_members_item.jsx';

import {Svgs} from '../../constants';

import PropTypes from 'prop-types';
import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';
import {searchPosts} from 'mattermost-redux/actions/search'
import {getChannel} from 'mattermost-redux/selectors/entities/channels';
import * as UserUtils from 'mattermost-redux/utils/user_utils';


// const dispatch = window.store.dispatch;
// const getState = window.store.getState;

export default class ChannelHeaderButton extends React.PureComponent {
    static propTypes = {
        channelId: PropTypes.string.isRequired,
        state: PropTypes.object.isRequired,
          channelName: PropTypes.string.isRequired,
        theme: PropTypes.object.isRequired,
        directChannels: PropTypes.array.isRequired,
        teamId: PropTypes.string.isRequired,
        channel: PropTypes.object.isRequired,
        actions: PropTypes.shape({
            startMeeting: PropTypes.func.isRequired,
            showRecordings: PropTypes.func.isRequired,
        }).isRequired
    }

    constructor(props) {
        super(props);

        this.state = {
            showPopover: false,
            rowStartHover: false,
            rowStartWithTopicHover: false,
            rowShareHover: false,
            showModal: false,
            shareModal: false,

        };
    }

    rowStartShowHover = () => {
        this.setState({rowStartHover: true});
    }

    rowStartHideHover = () => {
        this.setState({rowStartHover: false});
    }

    rowStartWithTopicShowHover = () => {
        this.setState({rowStartWithTopicHover: true});
    }

    rowStartWithTopicHideHover = () => {
        this.setState({rowStartWithTopicHover: false});
    }

    rowShareShowHover = () => {
        this.setState({rowShareHover: true});
    }

    rowShareHideHover = () => {
        this.setState({rowShareHover: false});
    }

    resetHover = () => {
        this.rowStartHideHover();
        this.rowStartWithTopicHideHover();
        this.rowShareHideHover();
    }

    showModal = () => {
        this.setState({showPopover: false, showModal: true, shareModal: false});
        this.resetHover();
    }

    searchRecordings =  () => {

        this.props.actions.showRecordings();
    }

    hideModal = () => {
        this.setState({showModal: false});
        this.resetHover();
    }
    hideRecordingsModal = () => {
        this.setState({showRecordingsModal: false});
        this.resetHover();
    }

    startMeeting = async () => {
      //console.log("all dm meetings" + JSON.stringify(this.props.directChannels));
      //const channel = getState().
      //console.log(JSON.stringify(this.props.channel));
        await this.props.actions.startMeeting(this.props.channelId, "",this.props.channel.display_name);
        this.setState({showPopover: false});
        this.resetHover();
    }




    render() {

      if (this.props.channelId === '') {
          return <div/>;
      }

      var channel = getChannel(this.props.state, this.props.channelId);
      var channelName = channel.display_name;


      // var fullname = UserUtils.getFullName(channel.teammate_id);
      console.log("aaaa "+ JSON.stringify(channel) );


        const style = getStyle(this.props.theme);

        let popoverButton = (
                <div
                    className='more-modal__button'
                >

                <a
                    className='btn  btn-link'

                    onClick={this.searchRecordings}
                >

                  {'View Recordings'}
                </a>

                </div>
            );



        return (

            <div>
                <div
                    id='bbbChannelHeaderPopover'
                    className={this.state.showPopover ? 'channel-header__icon active' : 'channel-header__icon'}
                >
                    <OverlayTrigger
                        trigger={['hover', 'focus']}
                        delayShow={400}
                        placement='bottom'
                        overlay={(
                            <Tooltip id='bbbChannelHeaderTooltip'>
                                {'BigBlueButton'}
                            </Tooltip>
                        )}
                    >
                        <div
                            id='bbbChannelHeaderButton'
                            onClick={(e) => {
                                this.setState({popoverTarget: e.target, showPopover: !this.state.showPopover});
                            }}
                        >
                            <span
                                style={style.iconStyle}
                                aria-hidden='true'
                                dangerouslySetInnerHTML={{__html: Svgs.SHARE}}
                            />
                        </div>
                    </OverlayTrigger>
                    <Overlay
                        rootClose={true}
                        show={this.state.showPopover}
                        target={() => this.state.popoverTarget}
                        onHide={() => this.setState({showPopover: false})}
                        placement='bottom'
                    >
                        <Popover
                            id='bbbPopover'
                            style={this.props.channel.type === "D" ?
                              style.popoverDM: style.popover }
                        >
                            <div style={this.props.channel.type === "D" ?
                              style.popoverBodyDM : style.popoverBody}>
                              {this.props.channel.type === "D" ?
                                <PopoverListMembersItem
                                  onItemClick = {this.startMeeting}
                                  cam = {1}
                                  text = {<span>{'Call '} <strong>{channelName}</strong></span>}
                                  theme = {this.props.theme}
                                /> :
                                <PopoverListMembersItem
                                  onItemClick = {this.startMeeting}
                                  cam = {1}
                                  text = {<span>{'Create a BigBlueButton Meeting'}</span>}
                                  theme = {this.props.theme}
                                /> }

                            </div>
                            {popoverButton}
                        </Popover>
                    </Overlay>
                </div>
        
            </div>
        );
    }
}

const getStyle = makeStyleFromTheme((theme) => {
    return {
        iconStyle: {
            position: 'relative',
            top: '-1px'
        },
        popover: {
            marginLeft: '-100px',
            maxWidth: '300px',
            height: '105px',
            width: '300px',
            background: theme.centerChannelBg
        },
        popoverBody: {
            maxHeight: '305px',
            overflow: 'auto',
            position: 'relative',
            width: '298px',
            left: '-14px',
            top: '-9px',
            borderBottom: '1px solid #D8D8D9'
        },
        popoverDM: {
            marginLeft: '-50px',
            maxWidth: '220px',
            height: '105px',
            width: '220px',
            background: theme.centerChannelBg
        },
        popoverBodyDM: {
            maxHeight: '305px',
            overflow: 'auto',
            position: 'relative',
            width: '218px',
            left: '-14px',
            top: '-9px',
            borderBottom: '1px solid #D8D8D9'
        },
    };
});
