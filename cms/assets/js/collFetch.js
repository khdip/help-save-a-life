document.addEventListener("DOMContentLoaded", function () {
    const tableBody = document.querySelector("#collTable tbody");
    const paginationContainer = document.querySelector(".pagination");
    let currentPage = 1;
    let currentSortBy = "";
    let currentOrder = "asc";

    function fetchData(page = 1, sortBy = "", order = "asc") {
        currentPage = page;
        currentSortBy = sortBy;
        currentOrder = order;

        const url = `/json/collection?page=${page}&SortBy=${sortBy}&Order=${order}`;
        fetch(url)
            .then(response => response.json())
            .then(data => {
                renderTable(data.CollHome);
                renderPagination(data.Paginator);
            })
            .catch(error => {
                console.error("Error fetching data:", error);
            });
    }

    function renderTable(items) {
        tableBody.innerHTML = "";
        items.forEach(item => {
            const row = document.createElement("tr");
            row.innerHTML = `
                <td>${item.Date}</td>
                <td>${item.AccountNumber}</td>
                <td>${item.AccountType}</td>
                <td>${item.Amount} ${item.Currency}</td>
            `;
            tableBody.appendChild(row);
        });
    }

    function renderPagination(paginator) {
        paginationContainer.innerHTML = "";

        if (paginator.Total === 0 || paginator.CountPaginate <= 1) return;

        const createPageItem = (page, label, isActive = false, isDisabled = false) => {
            const li = document.createElement("li");
            li.className = `page-item ${isActive ? "active" : ""} ${isDisabled ? "disabled" : ""}`;
            const a = document.createElement("a");
            a.className = `page-link ${isActive ? "bg-success text-light" : "text-success"}`;
            a.href = "#";
            a.textContent = label;
            a.addEventListener("click", (e) => {
                e.preventDefault();
                if (!isDisabled && page !== currentPage) {
                    fetchData(page, currentSortBy, currentOrder);
                }
            });
            li.appendChild(a);
            return li;
        };

        // Previous button
        paginationContainer.appendChild(
            createPageItem(currentPage - 1, "«", false, currentPage === 1)
        );

        // Page numbers
        paginator.Pages.forEach(p => {
            if (p.Order) {
                paginationContainer.appendChild(
                    createPageItem(p.Order, p.Order, p.Current)
                );
            } else {
                const li = document.createElement("li");
                li.className = "page-item disabled";
                li.innerHTML = `<a class="page-link text-success">...</a>`;
                paginationContainer.appendChild(li);
            }
        });

        // Next button
        paginationContainer.appendChild(
            createPageItem(currentPage + 1, "»", false, currentPage === paginator.CountPaginate)
        );
    }

    // Sort header click
    document.querySelectorAll(".sort").forEach(header => {
        header.addEventListener("click", function (e) {
            e.preventDefault();
            const sortBy = this.dataset.sort;
            const newOrder = (currentSortBy === sortBy && currentOrder === "asc") ? "desc" : "asc";
            fetchData(1, sortBy, newOrder);
        });
    });

    // Initial load
    fetchData();
});

