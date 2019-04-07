'use strict';

const e = React.createElement; 
class StoryTab extends React.Component {
    constructor(props) {
        super(props);
        this.state = { liked: false };
    }

    function clickedButton() {
        if (this.state.liked) {
            this.state.liked = false;
            return "You unliked this...";
        }
        this.state.liked = true;
        return "You liked this!";
    }

    render() {
        if (this.state.liked) {
            this.state.liked = false;
            return e(
                'button',
                { onClick: () => this.setState({ liked: true }) },
                'Like'
            ), 'You liked this.';
        }

        return e(
            'button',
            { onClick: () => this.setState({ liked: true }) },
            'Like'
        );
    }
}

const domContainer = document.querySelector('#view_story_tab');
ReactDOM.render(e(StoryTab), domContainer);
