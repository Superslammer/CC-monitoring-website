<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    
    <!-- Boostrap css -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">

    <!-- Custom css -->
    <link href="./sass/main.css" rel="stylesheet">

    <title>Dashboard</title>
</head>
<body>
    <div class="top">
        <div class="header">
            <?php include "./includes/header.php" ?>
        </div>
        <div class="container">
            <div class="row">
                <div class="col">
                    <div class="row contentBox energy">
                        <div class="col-4">
                            <p>Energy</p>
                            <p>Total energy: 1111 RF</p>
                            <p>Gain/Loss: +10 RF</p>
                            <br>
                            <select>
                                <option value="RF" selected>RF</option>
                                <option value="EU">EU</option>
                            </select>
                        </div>
                        <div class="col-8">
                            <p>Net gain over time</p>
                        </div>
                    </div>
                </div>
            </div>
            <div class="row">
                <div class="col contentBox">
                    <p>Bottom Text</p>
                </div>
            </div>
        </div>
    </div>
    <!-- Boostrap js and Popper -->
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p" crossorigin="anonymous"></script>
</body>
</html>