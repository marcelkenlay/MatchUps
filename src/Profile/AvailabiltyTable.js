import React from 'react';
import './Stylesheets/AvailabilityTable.css';
import UserProfile from './UserProfile';
import Button from "react-bootstrap/es/Button";

var axios = require('axios');

export default class AvailabiltyTable extends React.Component {
  state = {
    availability: [[0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                   [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                   [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                   [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                   [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                   [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                   [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]],
    isSubmitted: false,
    isMouseDown: false,
  }

  numToDay(i) {
    switch (i) {
      case 0:
        return "Monday";
      case 1:
        return "Tuesday";
      case 2:
        return "Wednesday";
      case 3:
        return "Thursday";
      case 4:
        return "Friday";
      case 5:
        return "Saturday";
      case 6:
        return "Sunday";
      default:
        return "Error undefined index of day ("+ i + ")";
    }
  }

  componentDidMount() {
    var username = UserProfile.getName();
    var _this = this;
    axios.get("/getuseravail?username=" + username)
         .then(function(response) {
           var bitmaps = new Array(7);
           for (var i = 0; i < 7; i++) {
             bitmaps[i] = parseInt(response.data[i]);
           }

           var avail = _this.bitmapToMatrix(bitmaps);

           _this.setState({
             availability: avail
           })
         });
  }

  createTable = () => {
    let rows = []
    let row = []
    row.push(<td></td>)

    // Add times
    for (let i = 8; i < 24; i++) {
      row.push(<td>{i}:00</td>)
    }
    rows.push(<tr>{row}</tr>)

    // Add each cell
    for (let i = 0; i < 7; i++) {
      row = []

      // Add day
      row.push(<td>{this.numToDay(i)}</td>)
      for (let j = 0; j < 16; j++) {
        row.push(<td
          onMouseDown={e => {
            this.state.isMouseDown = true;
            this.toggle(i, j, e);
            return false;
            }}
          onMouseOver={e => {
            if(this.state.isMouseDown)
              this.toggle(i, j, e)
          }}
          onMouseUp={e =>{
            this.state.isMouseDown = false;
          }}
          style = {{selectstart:'false'}} //prevents text selection while dragging
          id={"av" + i + "_" + j}
          class={"avBox " + (this.state.availability[i][j] ? "checked" : "unchecked")}
          ></td>)
      }

      // Add 'clearRow option'
      row.push(<td onClick={e => this.toggleDay(i, false)} id={"clrRw_" + i} class="toggleDay">clear day</td>)
      row.push(<td onClick={e => this.toggleDay(i, true)} id={"clrRw_" + i} class="toggleDay">fill day</td>)
      rows.push(<tr>{row}</tr>)
    }
    return rows
  }

  // Toggle selection of an availability slot
  toggle(i, j, e) {
    e.preventDefault();

    if (this.state.availability[i][j]) {
      // Deselect
      document.getElementById("av" + i + "_" + j).classList.remove("checked");
      document.getElementById("av" + i + "_" + j).classList.add("unchecked");
      this.state.availability[i][j] = 0;
    } else {
      // Select
      document.getElementById("av" + i + "_" + j).classList.remove("unchecked");
      document.getElementById("av" + i + "_" + j).classList.add("checked");
      this.state.availability[i][j] = 1;
    }

  }

  // Convert the state availability matrix to a bitmap for each day
  matrixToBitmap() {
    var av = this.state.availability;
    var res = new Array(7);

    for (var day = 0; day < 7; day++) {
      var bitmap = 0;

      for (var time = 0; time < 16; time++) {
        if (av[day][time]) {
          bitmap = bitmap | 0x1;
        }

        bitmap = bitmap << 1;
      }

      res[day] = bitmap >> 1;
    }

    return res;
  }

  // Convert bitmaps into the state availability matrix
  bitmapToMatrix(bitmaps) {
    // Create new matrix
    var matrix = new Array(7) // 7 rows for the days
    for (var row = 0; row < 7; row++) {
      matrix[row] = new Array(16)
    }

    // Process first bitmap (rows 0, 1, 2 and 3)
    // Read and shift
    for (var row = 6; row >= 0; row--) {
      var bitmap = bitmaps[row];

      for (var col = 15; col >= 0; col--) {
        matrix[row][col] = bitmap & 0x1;
        bitmap = bitmap >> 1;
      }
    }

    return matrix;
  }

  // Submit the modified availability to the server
  update(e) {
    e.preventDefault();
    var _this = this;

    // Create new availability matrix
    var newAvails = new Array(7) // 7 rows for the days

    for (var row = 0; row < 7; row++) {
      newAvails[row] = new Array(16)

      for (var col = 0; col < 16; col++) {
        var cell = document.getElementById("av" + row + "_" + col);
        if (cell.classList.contains("checked")) {
          newAvails[row][col] = 1;
        }
      }
    }

    // Set state to new matrix
    this.setState({
      availability: newAvails
    })

    // Convert matrix to two bitmap values
    var bitmaps = this.matrixToBitmap();
    var username = UserProfile.getName();

    // Send the bitmap values to the server
    axios.get("/updateuseravail?username=" + username +
              "&mon="  + bitmaps[0] +
              "&tues=" + bitmaps[1] +
              "&weds=" + bitmaps[2] +
              "&thurs=" + bitmaps[3] +
              "&fri="  + bitmaps[4] +
              "&sat="  + bitmaps[5] +
              "&sun="  + bitmaps[6])
         .then(function(response) {
           if (response.data == "fail\n") {
             alert("Failed to update availability, please try again.")
           } else {
             var tick = document.getElementById('tabletick');
             tick.innerHTML = "✓";
             setTimeout(function() {
               tick.innerHTML = "";
             }, 2500);
           }
         });
  }

  toggleDay(day, toOn) {
    for (var time = 0; time < 16; time++) {
      if (toOn) {
        // Select
        document.getElementById("av" + day + "_" + time).classList.remove("unchecked");
        document.getElementById("av" + day + "_" + time).classList.add("checked");
        this.state.availability[day][time] = 1;
      } else {
        // Deselect
        document.getElementById("av" + day + "_" + time).classList.remove("checked");
        document.getElementById("av" + day + "_" + time).classList.add("unchecked");
        this.state.availability[day][time] = 0;
      }
    }
  }

  render() {
    return (
      <div>
        <table>
          {this.createTable()}
        </table>
        <br />
        <p>Click a box to toggle availability between <span class='selectedtext checked'>selected</span> and <span class='unselectedtext unchecked'>unselected</span></p>

        <a><h3 class='flasher' onClick={e => this.update(e)}>Click to Save!<span id='tabletick'></span></h3></a>
      </div>
    )
  }
}
