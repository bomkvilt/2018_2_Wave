TRUNCATE apps CASCADE;
-- TRUNCATE userapp CASCADE;

INSERT INTO apps(link, name, image, about, category)
    VALUES (
        '/terminal',
        'Terminal',
        'img/app_covers/terminal.jpg',
        'Simple console, just like in your favorite Ubuntu',
        '2018_2'
    );

INSERT INTO apps(link, name, image, about, category)
    VALUES (
        '/snake',
        'Snake',
        'img/app_covers/snake.jpg',
        'Simple Snake game',
        '2018_2'
    );

INSERT INTO apps(link, url, name, image, about, category)
    VALUES (
        '/proto',
        'https://rasseki.pro/',
        'Quizzy',
        'img/app_covers/proto.jpg',
        'Simple duel quiz by __proto__',
        '2018_2'
    );

INSERT INTO apps(link, url, name, image, about, category)
    VALUES (
        '/blep',
        'https://blep.me',
        'Blep',
        'img/app_covers/blep.jpg',
        'Elegant game about imagination by Stacktivity',
        '2018_2'
    );

INSERT INTO apps(link, url, name, image, about, category)
    VALUES (
        '/catchthealien',
        'https://itberries-frontend.herokuapp.com/',
        'Catch the alien!',
        'img/app_covers/catch.jpg',
        'Dont let the alien leave the field! Game by ItBerries',
        '2018_1'
    );

INSERT INTO apps(link, url, name, image, about, category)
    VALUES (
        '/guardians',
        'https://chunk-frontend.herokuapp.com/',
        'Guardians',
        'img/app_covers/guardians.jpg',
        'Unusual 3D multi-player puzzle. Game by Chunk',
        '2017_2'
    );

INSERT INTO apps(link, url, name, image, about, category)
    VALUES (
        '/rhytmblast',
        'https://glitchless.surge.sh/',
        'Rhythm Blast',
        'img/app_covers/rhytmblast.jpg',
        'Space arcanoid on steroids. Game by Glitchless',
        '2017_2'
    );

INSERT INTO apps(link, url, name, image, about, category)
    VALUES (
        '/ketnipz',
        'https://playketnipz.ru/',
        'Ketnipz',
        'img/app_covers/ketnipz.jpg',
        'Arcade game with cartoony graphics by DeadMolesStudio',
        '2018_2'
    );

INSERT INTO apps(link, url, name, image, about, category)
    VALUES (
        '/kekmate',
        'https://kekmate.tech/',
        'Chessmate',
        'img/app_covers/chess.jpg',
        'Classic chess with advanced AI game by Parashutnaya Molitva',
        '2018_2'
    );

INSERT INTO apps(link, url, name, image, about, category)
    VALUES (
        '/rpsarena',
        'http://rpsarena.ru',
        'RPS Arena',
        'img/app_covers/rps.jpg',
        'Multiplayer version of classic Rock–Paper–Scissors game by 42',
        '2018_2'
    );
