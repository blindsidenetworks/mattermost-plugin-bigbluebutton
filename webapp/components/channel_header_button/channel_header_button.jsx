/*
Copyright 2018 Blindside Networks

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// const React = window.react;
import React from 'react';
const {Overlay, OverlayTrigger, Popover, Tooltip} = window['react-bootstrap'];

import PopoverListMembersItem from './popover_list_members_item.jsx';

import {Svgs} from '../../constants';

import PropTypes from 'prop-types';
import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';
import {searchPosts} from 'mattermost-redux/actions/search'
import {getChannel} from 'mattermost-redux/selectors/entities/channels';
import * as UserUtils from 'mattermost-redux/utils/user_utils';

export default class ChannelHeaderButton extends React.PureComponent {
  static propTypes = {
    channelId: PropTypes.string.isRequired,
    state: PropTypes.object.isRequired,
    channelName: PropTypes.string.isRequired,
    theme: PropTypes.object.isRequired,
    directChannels: PropTypes.array.isRequired,
    teamId: PropTypes.string.isRequired,
    channel: PropTypes.object.isRequired,
    actions: PropTypes.shape({startMeeting: PropTypes.func.isRequired, showRecordings: PropTypes.func.isRequired}).isRequired
  }

  constructor(props) {
    super(props);

    this.state = {
      showPopover: false
    };
  }

  searchRecordings = () => {
    this.props.actions.showRecordings();
  }

  startMeeting = async () => {
    await this.props.actions.startMeeting(this.props.channelId, "", this.props.channel.display_name);
    this.setState({showPopover: false});
  }

  render() {

    if (this.props.channelId === '') {
      return <div/>;
    }

    var channel = getChannel(this.props.state, this.props.channelId);
    var channelName = channel.display_name;

    const style = getStyle(this.props.theme);

    let popoverButton = (<div className='more-modal__button'>

      <a className='btn  btn-link' onClick={this.searchRecordings}>

        {'View Recordings'}
      </a>

    </div>);

    return (<div>
      <div id='bbbChannelHeaderPopover' className={this.state.showPopover
          ? 'channel-header__icon active'
          : 'channel-header__icon'}>
        <OverlayTrigger trigger={['hover', 'focus']} delayShow={400} placement='bottom' overlay={(<Tooltip id='bbbChannelHeaderTooltip'>
            {'BigBlueButton'}
          </Tooltip>)}>
          <div id='bbbChannelHeaderButton' onClick={(e) => {
              this.setState({
                popoverTarget: e.target,
                showPopover: !this.state.showPopover
              });
            }}>
            <span style={style.iconStyle} aria-hidden='true' dangerouslySetInnerHTML={{
                __html: Svgs.SHARE
              }}/>
          </div>
        </OverlayTrigger>
        <Overlay rootClose={true} show={this.state.showPopover} target={() => this.state.popoverTarget} onHide={() => this.setState({showPopover: false})} placement='bottom'>
          <Popover id='bbbPopover' style={this.props.channel.type === "D"
              ? style.popoverDM
              : style.popover}>
            <div style={this.props.channel.type === "D"
                ? style.popoverBodyDM
                : style.popoverBody}>
              {
                this.props.channel.type === "D"
                  ? <PopoverListMembersItem onItemClick={this.startMeeting} cam={1} text={<span> {
                        'Call '
                      }
                      <strong>{channelName}</strong>
                    </span>} theme={this.props.theme}/>
                  : <PopoverListMembersItem onItemClick={this.startMeeting} cam={1} text={<span> {
                        'Create a BigBlueButton Meeting'
                      }
                      </span>} theme={this.props.theme}/>
              }

            </div>
            {popoverButton}
          </Popover>
        </Overlay>
      </div>

    </div>);
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
    }
  };
});
