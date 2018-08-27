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

import PropTypes from 'prop-types';

import {Svgs} from '../../constants';
import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';

export default class PopoverListMembersItem extends React.PureComponent {
  static propTypes = {
    onItemClick: PropTypes.func.isRequired,
    text: PropTypes.element.isRequired,
    cam: PropTypes.number.isRequired,
    theme: PropTypes.object.isRequired
  };

  constructor(props) {
    super(props);
    this.state = {
      rowStartHover: false
    };
  }

  handleClick = () => {
    this.props.onItemClick();
  };

  rowStartShowHover = () => {
    this.setState({rowStartHover: true});
  }

  rowStartHideHover = () => {
    this.setState({rowStartHover: false});
  }

  render() {
    const style = getStyle(this.props.theme);

    return (<div onMouseEnter={this.rowStartShowHover} onMouseLeave={this.rowStartHideHover} onClick={this.handleClick} style={this.state.rowStartHover
        ? style.popoverRowHover
        : style.popoverRowNoHover}>
      <span style={style.popoverIcon} className='pull-left' dangerouslySetInnerHTML={this.props.cam == 1
          ? {
            __html: Svgs.BBBCAM
          }
          : {
            __html: Svgs.VID_CAM_PLAY
          }} aria-hidden='true'/>
      <div style={style.popoverRow}>
        <div style={style.popoverText}>
          {this.props.text}
        </div>
      </div>
    </div>)
  }

}

const getStyle = makeStyleFromTheme((theme) => {
  return {

    popoverRow: {
      border: 'none',
      cursor: 'pointer',
      height: '50px',
      margin: '1px 0',
      overflow: 'auto',
      padding: '6px 19px 0 10px'
    },
    popoverRowNoHover: {
      borderLeft: '3px solid',
      borderColor: theme.centerChannelBg,
      fontWeight: 'normal'
    },
    popoverRowHover: {
      borderLeft: '3px solid transparent',
      borderColor: theme.linkColor,
      background: changeOpacity(theme.linkColor, 0.08)
    },
    popoverText: {
      fontWeight: 'inherit',
      fontSize: '14px',
      position: 'relative',
      top: '10px',
      left: '4px'
    },
    popoverIcon: {
      margin: '0',
      paddingLeft: '16px',
      position: 'relative',
      top: '12px',
      fontSize: '20px',
      fill: theme.buttonBg
    }
  };
});
