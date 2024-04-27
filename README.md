# TrialTanksUnity3D
Tugas Besar 2 IF2211 Strategi Algoritma 2023/2024

Pemanfaatan Algoritma IDS dan BFS dalam Permainan WikiRace

Anggota Kelompok:
10023500 - Miftahul Jannah
13518108 - Vincent Hasiholan
13522106 - Christopher Brian

## Algoritma Iterative Deepening Search (IDS)
Dalam Tugas Besar ini, algoritma IDS yang kamu gunakan berupa Depth Limited Search (DLS) secara iteratif dengan depth limit yang meningkat di setiap iterasi. DLS akan mencari semua kemungkinan path dari page yang sedang dicek ke page yang bisa dicapai menggunakan link secara rekursif. Jika target page ditemukan, maka path ke page tersebut akan dikembalikan. Jika tidak, depth limit akan ditingkatkan pada iterasi berikutnya, dan pencarian akan terus dilakukan hingga target page ditemukan pada path atau depth limit maksimum dicapai, yang kami atur menjadi maksimal 10 degree. Algoritma ini menggabunkan completeness dari Depth First Search (DFS) dan efisiensi dari Breadth First Search (BFS).

## Algoritma Breadth First Search (BFS)
Kami juga menggunakan algoritma BFS sebagai pembanding. Algoritma yang kami gunakan berupa BFS direksional, sehingga proses pencarian dimulai dari source page dan juga target page, dengan tujuan untuk mencari page di tengah-tengah yang bisa dicapai dari kedua page tersebut. Terdapat dua daftar page, yaitu unvisitedForward, yang berupa page yang belum dikunjungi pada arah maju (dari source page menuju target page), dan unvisitedBackward (dari arah sebaliknya). Algoritma ini akan mengembangkan scope pencarian secara iteratif dari kedua arah sampai menemukan page yang bisa diakses dari source dan target page. Kemudian, algoritma akan mencari semua path yang mungkin dari source page ke target page melalui page yang bisa diakses dari kedua arah tersebut.

## Requirement
PC/laptop sudah terinstall go v1.22.2 atau ke atas

## Langkah Penggunaan
1. Mengetik perintah go run wikiread.go pathfinding.go pada cmd atau terminal pada direktori yang memiliki file wikiread.go
2. Setelah itu, buka web browser (dalam kasus ini, google chrome), ketik localhost:9999 dan website muncul
3. Mengisi titik mulai pencarian dan titik tujuan dari pencarian
4. Menekan tombol submit ketika sudah yakin apa yang mau dicari
5. Aplikasi memproses inputan yang sudah disubmit dan menunjukkan hasilnya.