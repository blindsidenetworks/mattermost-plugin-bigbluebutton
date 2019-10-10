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


import React from 'react';
const {Overlay, OverlayTrigger, Popover, Tooltip} = window.ReactBootstrap;

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
    visible: PropTypes.bool.isRequired,
    actions: PropTypes.shape({startMeeting: PropTypes.func.isRequired, showRecordings: PropTypes.func.isRequired, closePopover: PropTypes.func.isRequired}).isRequired,
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
    this.close_the_popover()
  }
  close_the_popover = () =>{
    this.props.actions.closePopover();
    this.setState({showPopover: false});
  }

  render() {

    if (this.props.channelId === '') {
      return <div/>;
    }

    var channel = getChannel(this.props.state, this.props.channelId);
    var channelName = channel.display_name;

    const style = getStyle(this.props.theme);


    return (<div>
      <div >
        <OverlayTrigger trigger={['hover']} delayShow={400} ref = "overlay" placement='bottom' overlay={(<Tooltip id='bbbChannelHeaderTooltip'>
            {'BigBlueButton'}
          </Tooltip>)}>
          <div id='bbbChannelHeaderButton' onClick={(e) => {
              this.refs.overlay.hide();
              this.setState({
                popoverTarget: e.target,
                showPopover: !this.props.visible
              });
            }}
               style={style.foo}
          >
            <span style={style.iconStyle} aria-hidden='true' dangerouslySetInnerHTML={{
                __html: Svgs.SHARE
              }}/>
          </div>
        </OverlayTrigger>
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
    },
    foo: {
      width: '28px',
      height: '28px',
    }
  };
});
