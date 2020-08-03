 {{define "articleList"}}
 {{range $key,$article := .}}
    <div class="list-group-item list-group-item-action">
        <div class="media">
            <a href="/user/{{$article.UserName}}"><img src="{{$article.AvatarUrl}}" class="mr-3 rounded-circle" style="width:45px;height:45px; " alt="..."></a>
            <div class="media-body">
                <h5 class="mt-0 mb-1"><a class="text-dark" href="/user/{{$article.UserName}}/article/{{$article.Id.Hex}}" target="_blank">{{$article.Title}}</a></h5>
                <span class="badge badge-success">{{$article.Type}}</span>&nbsp;&nbsp;<span class="badge badge-secondary">{{$article.Created | fdate}}</span>
            </div>
        </div>
    </div>
{{end}}
{{end}}