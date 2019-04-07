'use strict';

const e = React.createElement; 
class StoryTab extends React.Component {
    constructor(props) {
        super(props);
        this.state = { hidden: false,
                       content: [],
                       title: '',
                       story: '',
                    
                        };

    }

    render() {
        if (this.state.hidden) {
            return (
                <h4>
        }
        return e(
            'button',
            { type: "button",
              onClick: () => this.setState({ liked: true }) ,
              class: "btn btn-secondary"},
            'Like'
        );
    }
}

const domContainer = document.querySelector('#view_story_tab');
ReactDOM.render(e(StoryTab), domContainer);
