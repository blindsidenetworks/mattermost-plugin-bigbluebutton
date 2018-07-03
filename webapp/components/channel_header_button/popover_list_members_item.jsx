const React = window.react;
import PropTypes from 'prop-types';

import {Svgs} from '../../constants';
import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';
import {FormattedHTMLMessage, FormattedMessage} from 'react-intl';

export default class PopoverListMembersItem extends React.PureComponent {
  static propTypes = {
        onItemClick: PropTypes.func.isRequired,
        text: PropTypes.element.isRequired,
        cam: PropTypes.number.isRequired,
        theme: PropTypes.object.isRequired,
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
    }

    rowStartHideHover = () => {
        this.setState({rowStartHover: false});
    }
    resetHover = () => {
        this.rowStartHideHover();
    }

    render(){
      const style = getStyle(this.props.theme);

        return (
          <div
              onMouseEnter={this.rowStartShowHover}
              onMouseLeave={this.rowStartHideHover}
              onClick={this.handleClick}
              style={this.state.rowStartHover ?
                style.popoverRowHover : style.popoverRowNoHover}
          >
              <span
                  style={style.popoverIcon}
                  className='pull-left'
                  dangerouslySetInnerHTML={this.props.cam == 1 ? {__html: Svgs.BBBCAM} : {__html: Svgs.VID_CAM_PLAY}}
                  aria-hidden='true'
              />
              <div style={style.popoverRow}>
                  <div style={style.popoverText}>
                      {this.props.text}
                  </div>
              </div>
          </div>
        )
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
            background: changeOpacity(theme.linkColor, 0.08),

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
