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
    icon: PropTypes.string.isRequired,
    theme: PropTypes.object.isRequired,
    ariaLabel: PropTypes.string,
  };

  constructor(props) {
    super(props);
    this.state = {
      rowStartHover: false,
    };
  }

  handleClick = () => {
    this.props.onItemClick();
  };

  rowStartShowHover = () => {
    this.setState({rowStartHover: true});
  };

  rowStartHideHover = () => {
    this.setState({rowStartHover: false});
  };

  render() {
    const style = getStyle(this.props.theme);

    console.log(this.props);
    console.log(this.props.icon);
    console.log(Svgs[this.props.icon]);

    return (
      <button
        aria-label={this.props.ariaLabel}
        className={'style--none'}
        onMouseEnter={this.rowStartShowHover}
        onMouseLeave={this.rowStartHideHover}
        onClick={this.handleClick}
        style={this.state.rowStartHover ? {...style.popoverRowBase, ...style.popoverRowHover} : style.popoverRowBase}
      >
        <span
          style={style.popoverIcon}
          className='pull-left'
          // dangerouslySetInnerHTML={this.props.cam == 1 ? {__html: Svgs.BBBCAM} : {__html: Svgs.VID_CAM_PLAY}}
          dangerouslySetInnerHTML={{__html: Svgs[this.props.icon]}}
          aria-hidden='true'
        />
        <div style={style.popoverRow}>
          <div style={style.popoverText}>
            {this.props.text}
          </div>
        </div>
      </button>);
  }
}

const getStyle = makeStyleFromTheme((theme) => {
  return {

    popoverRow: {
      padding: '0 12px',
    },
    popoverRowBase: {
      borderLeft: '3px solid',
      borderColor: theme.centerChannelBg,
      fontWeight: 'normal',
      width: '100%',
      display: 'flex',
      height: '60px',
      padding: '6px 12px',
    },
    popoverRowHover: {
      borderLeft: '3px solid transparent',
      borderColor: theme.linkColor,
      background: changeOpacity(theme.linkColor, 0.08),
      width: '100%',
    },
    popoverText: {
      fontWeight: 'inherit',
      fontSize: '13px',
      textAlign: 'left',
    },
    popoverIcon: {
      height: '90%',
      padding: '4px',
    },
  };
});
